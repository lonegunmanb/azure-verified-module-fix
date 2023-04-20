package pkg

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type HclFile struct {
	*hcl.File
	WriteFile *hclwrite.File
}

func ParseConfig(config []byte, filename string) (*HclFile, hcl.Diagnostics) {
	file, rDiag := hclsyntax.ParseConfig(config, filename, hcl.InitialPos)
	writeFile, wDiag := hclwrite.ParseConfig(config, filename, hcl.InitialPos)
	if rDiag.HasErrors() || wDiag.HasErrors() {
		return nil, rDiag.Extend(wDiag)
	}
	return &HclFile{file, writeFile}, hcl.Diagnostics{}
}

func (f *HclFile) GetBlock(i int) *HclBlock {
	block := f.Body.(*hclsyntax.Body).Blocks[i]
	writeBlock := f.WriteFile.Body().Blocks()[i]
	return NewHclBlock(block, writeBlock)
}

func (f *HclFile) AutoFix() {
	for i, b := range f.Body.(*hclsyntax.Body).Blocks {
		hclBlock := f.GetBlock(i)
		if b.Type == "resource" || b.Type == "data" {
			resourceBlock := BuildResourceBlock(hclBlock, f.File)
			resourceBlock.AutoFix()
		} else if b.Type == "locals" {
			localsBlock := BuildLocalsBlock(hclBlock, f.File)
			localsBlock.AutoFix()
		}
	}
}