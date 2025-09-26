package commercial

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewCorporation(t *testing.T) {
	corp := NewCorporation(IsHumanPlayable(true))
	corp.AddInvestor()
	corp.StartIndustryLines(WorkforceManagement, 1)
	fmt.Println("start wealth", corp.Wealth)
	corp.WorkOnIndustries()
	// report := strings.Join(corp.RecordSheet(), "\n")

	// fmt.Println(report)
	// fmt.Println(corp.profitPool)
	fmt.Println(corp.MegaTrade(5))
	corp.WorkOnIndustries()
	corp.CalculateProfit()
	corp.PayStaff()
	fmt.Println(corp.ProfitShare(1))
	report := strings.Join(corp.RecordSheet(), "\n")

	fmt.Println(report)
	if err := ActionProtectCompanyThrueeLegalRedTape.AssesPrequisites(corp); err == nil {
		fmt.Println(corp.Commence(ActionProtectCompanyThrueeLegalRedTape))
	} else {
		fmt.Println(err)
	}
	fmt.Println(corp.NextDefenceDM)
}
