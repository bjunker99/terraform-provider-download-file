package download_file

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataDownloadFile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataDownloadFileRead,
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"output_file": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"output_base64sha256": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "Base64 Encoded SHA256 checksum of output file",
			},
			"output_md5": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "MD5 of output file",
			},
                        "output_sha": {
                                Type:        schema.TypeString,
                                Computed:    true,
                                ForceNew:    true,
                                Description: "SHA1 checksum of output file",
                        },
                        "output_sha256": {
                                Type:        schema.TypeString,
                                Computed:    true,
                                ForceNew:    true,
                                Description: "SHA256 checksum of output file",
                        },
			"output_size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				ForceNew: true,
			},
                        "verify_sha256": {
                                Type:        schema.TypeString,
                                Optional:    true,
                                Description: "SHA256 checksum to verify",
                        },
                        "verify_sha": {
                                Type:        schema.TypeString,
                                Optional:    true,
                                Description: "SHA checksum to verify",
                        },
                        "verify_md5": {
                                Type:        schema.TypeString,
                                Optional:    true,
                                Description: "MD5 checksum to verify",
                        },
		},
	}
}

func dataDownloadFileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	url := d.Get("url").(string)
	outputFile := d.Get("output_file").(string)

	err := DownloadFile(outputFile, url)

	if err != nil {
		return diag.FromErr(err)
	}

	fi, err := os.Stat(outputFile)
	if err != nil {
		return diag.FromErr(err)
	}

	sha1, sha256, base64sha256, md5, err := genFileShas(outputFile)

	if err != nil {
		return diag.FromErr(err)
	}

	err = Verify(d, sha256, sha1, md5)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("output_sha", sha1)
	d.Set("output_sha256", sha256)
	d.Set("output_base64sha256", base64sha256)
	d.Set("output_md5", md5)
	d.Set("output_size", fi.Size())

	d.SetId(url)

	return diags
}

func Verify(d *schema.ResourceData, sha256 string, sha string, md5 string) error {
	if v, ok := d.GetOk("verify_sha256"); ok {
		if v.(string) != sha256 {
			return errors.New("SHA256 signature mismatch")
		}
	}

        if v, ok := d.GetOk("verify_sha"); ok {
                if v.(string) != sha {
                        return errors.New("SHA signature mismatch")
                }
        }

        if v, ok := d.GetOk("verify_md5"); ok {
                if v.(string) != md5 {
                        return errors.New("MD5 signature mismatch")
                }
        }

	return nil
}

func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func genFileShas(filename string) (string, string, string, string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", "", "", "", fmt.Errorf("could not compute file '%s' checksum: %s", filename, err)
	}
	h := sha1.New()
	h.Write([]byte(data))
	sha1 := hex.EncodeToString(h.Sum(nil))

	h256 := sha256.New()
	h256.Write([]byte(data))
	shaSum := h256.Sum(nil)
	sha256 := hex.EncodeToString(h256.Sum(nil))
	sha256base64 := base64.StdEncoding.EncodeToString(shaSum[:])

	md5 := md5.New()
	md5.Write([]byte(data))
	md5Sum := hex.EncodeToString(md5.Sum(nil))

	return sha1, sha256, sha256base64, md5Sum, nil
}
