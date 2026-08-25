package importer

import "testing"

func TestParseAndValidatePayload(t *testing.T) {
	bundle, err := ParseJSON([]byte(`{"records":[{"title":"祖屋","narrative":"祖屋故事至少八字。","author":"研究员","amount":12,"attachments":[{"name":"来源.txt","media_type":"text/plain","checksum":"abc","size":5}]}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateBundle(bundle); err != nil {
		t.Fatal(err)
	}
	if bundle.Records[0].Title != "祖屋" || len(bundle.Records[0].Attachments) != 1 {
		t.Fatalf("bundle=%+v", bundle)
	}
}

func TestValidateRejectsDuplicates(t *testing.T) {
	payload := Payload{Records: []RecordInput{{Title: "重复", Narrative: "这是重复故事一。", Author: "a"}, {Title: "重复", Narrative: "这是重复故事二。", Author: "b"}}}
	bundle, err := ParsePayload(payload)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateBundle(bundle); err == nil {
		t.Fatal("expected duplicate error")
	}
}

func TestParseCSV(t *testing.T) {
	bundle, err := ParseCSV("title,narrative,author,amount,tags\n祖屋,祖屋来源故事文本,研究员,19,家乡|旧宅\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(bundle.Records) != 1 || len(bundle.Records[0].Tags) != 2 {
		t.Fatalf("bundle=%+v", bundle)
	}
}
