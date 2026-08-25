package application

import "genealogy-story-organizer/internal/importer"

func importerPayload() importer.Payload {
	return importer.Payload{Records: []importer.RecordInput{
		{Title: "祖母口述", Narrative: "祖母记得河畔旧居的故事。", Author: "研究员", Amount: 130, Tags: []string{"口述"}, Attachments: []importer.AttachmentInput{{Name: "口述.txt", MediaType: "text/plain", Checksum: "one", Size: 12}}},
		{Title: "族谱修订", Narrative: "族谱修订记录说明了支系变化。", Author: "研究员", Amount: 200, Tags: []string{"族谱"}, Attachments: []importer.AttachmentInput{{Name: "修订.pdf", MediaType: "application/pdf", Checksum: "two", Size: 20}}},
	}}
}
