package cmd

import "github.com/gwleclerc/adr/records"

// markSuperseded flags each target record as superseded by bySuperseder and
// back-links it. Unknown target IDs produce a warning instead of being ignored.
func markSuperseded(service *records.Service, bySuperseder string, targetIDs []string) {
	for _, id := range targetIDs {
		rcd, ok := service.GetRecord(id)
		if !ok {
			printWarning("superseded record %q not found, skipping", id)
			continue
		}
		rcd.Status = records.SUPERSEDED
		rcd.Superseders.Append(bySuperseder)
		if err := service.UpdateRecord(rcd); err != nil {
			printWarning("unable to update superseded record %q: %v", rcd.ID, err)
		}
	}
}

// warnUnknownRecords warns for every id that does not match an existing record.
func warnUnknownRecords(service *records.Service, ids []string) {
	for _, id := range ids {
		if _, ok := service.GetRecord(id); !ok {
			printWarning("superseder %q does not match any existing record", id)
		}
	}
}
