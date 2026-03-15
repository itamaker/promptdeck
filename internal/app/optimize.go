package app

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

func loadScores(path string) ([]CandidateScore, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read scores: %w", err)
	}

	var scores []CandidateScore
	if err := json.Unmarshal(body, &scores); err != nil {
		return nil, fmt.Errorf("decode scores: %w", err)
	}
	return scores, nil
}

func optimize(candidates []PromptCandidate, scores []CandidateScore, topN int) OptimizationReport {
	report := OptimizationReport{
		Candidates: len(candidates),
	}
	if topN <= 0 {
		topN = 3
	}

	scoreIndex := map[string]CandidateScore{}
	for _, score := range scores {
		scoreIndex[score.Index] = score
	}

	var ranked []RankedCandidate
	type bucket struct {
		count int
		sum   float64
	}
	effectIndex := map[string]*bucket{}

	for _, candidate := range candidates {
		score, ok := scoreIndex[candidate.Index]
		if !ok {
			continue
		}
		report.Scored++
		ranked = append(ranked, RankedCandidate{
			Index:  candidate.Index,
			Score:  score.Score,
			Vars:   candidate.Vars,
			Prompt: candidate.Prompt,
			Notes:  score.Notes,
		})
		for key, value := range candidate.Vars {
			if key == "index" {
				continue
			}
			bucketKey := key + "=" + fmt.Sprint(value)
			entry := effectIndex[bucketKey]
			if entry == nil {
				entry = &bucket{}
				effectIndex[bucketKey] = entry
			}
			entry.count++
			entry.sum += score.Score
		}
	}

	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].Score == ranked[j].Score {
			return ranked[i].Index < ranked[j].Index
		}
		return ranked[i].Score > ranked[j].Score
	})
	if len(ranked) > 0 {
		best := ranked[0]
		report.Best = &best
	}
	if len(ranked) > topN {
		report.Top = ranked[:topN]
	} else {
		report.Top = ranked
	}

	for key, entry := range effectIndex {
		parts := strings.SplitN(key, "=", 2)
		report.FactorEffects = append(report.FactorEffects, FactorEffect{
			Variable: parts[0],
			Value:    parts[1],
			Count:    entry.count,
			AvgScore: entry.sum / float64(entry.count),
		})
	}
	sort.Slice(report.FactorEffects, func(i, j int) bool {
		if report.FactorEffects[i].AvgScore == report.FactorEffects[j].AvgScore {
			if report.FactorEffects[i].Variable == report.FactorEffects[j].Variable {
				return report.FactorEffects[i].Value < report.FactorEffects[j].Value
			}
			return report.FactorEffects[i].Variable < report.FactorEffects[j].Variable
		}
		return report.FactorEffects[i].AvgScore > report.FactorEffects[j].AvgScore
	})
	if len(report.FactorEffects) > 8 {
		report.FactorEffects = report.FactorEffects[:8]
	}

	return report
}

func printOptimizationReport(report OptimizationReport) {
	fmt.Printf("Candidates: %d\n", report.Candidates)
	fmt.Printf("Scored: %d\n", report.Scored)
	if report.Best != nil {
		fmt.Printf("Best prompt: #%s score=%.3f\n", report.Best.Index, report.Best.Score)
	}
	if len(report.Top) > 0 {
		fmt.Println("Top prompts:")
		for _, candidate := range report.Top {
			fmt.Printf("- #%s score=%.3f vars=%v\n", candidate.Index, candidate.Score, candidate.Vars)
		}
	}
	if len(report.FactorEffects) > 0 {
		fmt.Println("Best factors:")
		for _, effect := range report.FactorEffects {
			fmt.Printf("- %s=%s avg=%.3f (%d)\n", effect.Variable, effect.Value, effect.AvgScore, effect.Count)
		}
	}
}
