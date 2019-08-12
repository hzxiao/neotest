package neo

import "github.com/CityOfZion/neo-go/pkg/core/transaction"

var (
	attrLookup = map[string]transaction.AttrUsage{
		"ContractHash":   transaction.ContractHash,
		"ECDH02":         transaction.ECDH02,
		"ECDH03":         transaction.ECDH03,
		"Script":         transaction.Script,
		"Vote":           transaction.Vote,
		"CertURL":        transaction.CertURL,
		"DescriptionURL": transaction.DescriptionURL,
		"Description":    transaction.Description,

		"Hash1":  transaction.Hash1,
		"Hash2":  transaction.Hash2,
		"Hash3":  transaction.Hash3,
		"Hash4":  transaction.Hash4,
		"Hash5":  transaction.Hash5,
		"Hash6":  transaction.Hash6,
		"Hash7":  transaction.Hash7,
		"Hash8":  transaction.Hash8,
		"Hash9":  transaction.Hash9,
		"Hash10": transaction.Hash10,
		"Hash11": transaction.Hash11,
		"Hash12": transaction.Hash12,
		"Hash13": transaction.Hash13,
		"Hash14": transaction.Hash14,
		"Hash15": transaction.Hash15,

		"Remark":   transaction.Remark,
		"Remark1":  transaction.Remark1,
		"Remark2":  transaction.Remark2,
		"Remark3":  transaction.Remark3,
		"Remark4":  transaction.Remark4,
		"Remark5":  transaction.Remark5,
		"Remark6":  transaction.Remark6,
		"Remark7":  transaction.Remark7,
		"Remark8":  transaction.Remark8,
		"Remark9":  transaction.Remark9,
		"Remark10": transaction.Remark10,
		"Remark11": transaction.Remark11,
		"Remark12": transaction.Remark12,
		"Remark13": transaction.Remark13,
		"Remark14": transaction.Remark14,
		"Remark15": transaction.Remark15,
	}
)

func ValidAttrUsage(usage string) bool {
	_, ok := attrLookup[usage]
	return ok
}
