package valueobjects

type CovidStats struct {
	CumulativeDeaths int
	NewDeaths        int
	CumulativeCases  int
	NewCases         int
}

func NewCovidStats(cumulativeDeaths, newDeaths, cumulativeCases, newCases int) CovidStats {
	return CovidStats{
		CumulativeDeaths: cumulativeDeaths,
		NewDeaths:        newDeaths,
		CumulativeCases:  cumulativeCases,
		NewCases:         newCases,
	}
}
