package character

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

func TestAbilityDescriptions(t *testing.T) {
	actual := AbilityDescriptions()
	assert.Equal(t, 6, len(actual))
}

func TestAbilityScoreModifier(t *testing.T) {
	actual := AbilityScoreModifier()
	assert.Equal(t, 30, len(actual))
	assert.Equal(t, 0, actual[0])
	assert.Equal(t, 10, actual[30])
}

func TestAbilityAssign(t *testing.T) {
	actual := AbilityAssign()
	actualKeys := GetAbilityRollingOptions()
	assert.Equal(t, 8, len(actual))
	assert.Equal(t, 8, len(actualKeys))
}

func TestAbilityScorePointCost(t *testing.T) {
	actual := AbilityScorePointCost()
	assert.Equal(t, 8, len(actual))
	assert.Equal(t, 2, actual[10])
}

func TestSkillAbilityLookup(t *testing.T) {
	actual := SkillAbilityLookup()
	assert.Equal(t, 18, len(actual))
	assert.Equal(t, "intelligence", actual["history"])
}

func TestAbilityArrayTemplate(t *testing.T) {
	actual := AbilityArrayTemplate()
	assert.Equal(t, 6, len(actual))
	assert.Equal(t, 0, actual["charisma"])
}

func TestGetPreGeneratedBaseAbilityArray(t *testing.T) {
	actual, sortOrder := GetPreGeneratedBaseAbilityArray([]int{18, 17, 16, 15, 14, 13})
	assert.Equal(t, 6, len(actual))
	expected := map[string]int{
		"strength":     18,
		"dexterity":    17,
		"constitution": 16,
		"intelligence": 15,
		"wisdom":       14,
		"charisma":     13,
	}
	assert.Equal(t, expected, actual)
	expectedSortOrder := []string{
		"strength",
		"dexterity",
		"constitution",
		"intelligence",
		"wisdom",
		"charisma",
	}
	assert.Equal(t, expectedSortOrder, sortOrder)
}

func TestGetBaseAbilityArray(t *testing.T) {
	// Given
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLoggerSugared := zap.New(observedZapCore).Sugar()
	sortOrder := []string{"dexterity", "constitution", "strength",
		"charisma", "wisdom", "intelligence"}
	rollingOption := "standard"

	// When
	actual, r, err := GetBaseAbilityArray(sortOrder, rollingOption, observedLoggerSugared)

	// Then
	assert.Equal(t, nil, err)
	assert.Equal(t, []int{15, 14, 13, 12, 10, 8}, r)
	expected := map[string]int{
		"strength":     13,
		"dexterity":    15,
		"constitution": 14,
		"intelligence": 8,
		"wisdom":       10,
		"charisma":     12,
	}
	assert.Equal(t, expected, actual)
	require.Equal(t, 1, observedLogs.Len())
}

func TestGetBaseAbilityArrayWithRolls(t *testing.T) {
	// Given
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLoggerSugared := zap.New(observedZapCore).Sugar()
	sortOrder := []string{"dexterity", "constitution", "strength",
		"charisma", "wisdom", "intelligence"}
	rollingOption := "common"

	// When
	actual, r, err := GetBaseAbilityArray(sortOrder, rollingOption, observedLoggerSugared)

	// Then
	assert.Equal(t, nil, err)
	assert.Equal(t, 6, len(r))
	assert.GreaterOrEqual(t, actual["dexterity"], actual["constitution"])
	assert.GreaterOrEqual(t, actual["constitution"], actual["strength"])
	assert.GreaterOrEqual(t, actual["strength"], actual["charisma"])
	assert.GreaterOrEqual(t, actual["charisma"], actual["wisdom"])
	assert.GreaterOrEqual(t, actual["wisdom"], actual["intelligence"])
	require.Equal(t, 7, observedLogs.Len()) // 6 dice rolls and the sorted map
}

func TestGetPreGeneratedAbilityArray(t *testing.T) {
	Raw := []int{18, 17, 16, 15, 14, 13}
	ArchetypeBonus := AbilityArrayTemplate()
	ArchetypeBonus["charisma"] = 2
	ArchetypeBonus["intelligence"] = 1
	ArchetypeBonusIgnored := false
	LevelChangeIncrease := AbilityArrayTemplate()
	LevelChangeIncrease["dexterity"] = 2
	AdditionalBonus := AbilityArrayTemplate()
	AdditionalBonus["strength"] = 2
	ctxRef := "TestGetPreGeneratedAbilityArray"
	isMonsterOrGod := false
	a := GetPreGeneratedAbilityArray(Raw, ArchetypeBonus,
		ArchetypeBonusIgnored, LevelChangeIncrease,
		AdditionalBonus, ctxRef, isMonsterOrGod)
	fmt.Println(a.ToPrettyString())
	assert.Equal(t, 20, a.Values["strength"])
	assert.Equal(t, 15, a.Values["charisma"])
	assert.Equal(t, 19, a.Values["dexterity"])
	assert.Equal(t, 16, a.Values["intelligence"])
	actual, _ := a.GetModifier("strength")
	assert.Equal(t, 5, actual)
	actual, _ = a.GetModifier("intelligence")
	assert.Equal(t, 3, actual)
}

func TestGetAbilityArray(t *testing.T) {
	// Given
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLoggerSugared := zap.New(observedZapCore).Sugar()
	rollingOption := "standard"
	sortOrder := []string{"dexterity", "constitution", "strength",
		"charisma", "wisdom", "intelligence"}
	ArchetypeBonus := AbilityArrayTemplate()
	ArchetypeBonus["charisma"] = 2
	ArchetypeBonus["intelligence"] = 1
	ArchetypeBonusIgnored := true
	LevelChangeIncrease := AbilityArrayTemplate()
	LevelChangeIncrease["dexterity"] = 2
	AdditionalBonus := AbilityArrayTemplate()
	AdditionalBonus["strength"] = 2
	ctxRef := "TestGetAbilityArray"
	isMonsterOrGod := true

	// When
	a, err := GetAbilityArray(rollingOption, sortOrder, ArchetypeBonus,
		ArchetypeBonusIgnored, LevelChangeIncrease, AdditionalBonus,
		ctxRef, isMonsterOrGod, observedLoggerSugared)

	// Then
	assert.Equal(t, nil, err)
	fmt.Println(a.ToPrettyString())
	assert.Equal(t, 15, a.Values["strength"])
	assert.Equal(t, 12, a.Values["charisma"])
	assert.Equal(t, 17, a.Values["dexterity"])
	assert.Equal(t, 8, a.Values["intelligence"])
	actual, _ := a.GetModifier("strength")
	assert.Equal(t, 2, actual)
	actual, _ = a.GetModifier("intelligence")
	assert.Equal(t, -1, actual)
	allLogs := observedLogs.All()

	ctxMap := allLogs[0].ContextMap()
	tStr, _ := ctxMap["rawValues"].(string)
	assert.Equal(t, "[15, 14, 13, 12, 10, 8]", tStr)
	tStr, _ = ctxMap["sortedValues"].(string)
	assert.Equal(t, "{\"Strength\": 13, \"Dexterity\": 15, \"Constitution\": 14, "+
		"\"Intelligence\":  8, \"Wisdom\": 10, \"Charisma\": 12}", tStr)

	assert.Equal(t, "GetAbilityArray", allLogs[len(allLogs)-1].Message)
}

func TestAdjustValues(t *testing.T) {
	// Given
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLoggerSugared := zap.New(observedZapCore).Sugar()
	rollingOption := "standard"
	sortOrder := []string{"dexterity", "constitution", "strength",
		"charisma", "wisdom", "intelligence"}
	ArchetypeBonus := AbilityArrayTemplate()
	ArchetypeBonusIgnored := false
	LevelChangeIncrease := AbilityArrayTemplate()
	AdditionalBonus := AbilityArrayTemplate()
	ctxRef := "TestAdjustValues"
	isMonsterOrGod := false

	a, err := GetAbilityArray(rollingOption, sortOrder, ArchetypeBonus,
		ArchetypeBonusIgnored, LevelChangeIncrease, AdditionalBonus,
		ctxRef, isMonsterOrGod, observedLoggerSugared)
	assert.Equal(t, nil, err)
	a.AdjustValues("ArchetypeBonus", "charisma",
		2, observedLoggerSugared)
	a.AdjustValues("ArchetypeBonus", "intelligence",
		1, observedLoggerSugared)
	assert.Equal(t, 14, a.Values["charisma"])
	assert.Equal(t, 9, a.Values["intelligence"])
	a.AdjustValues("LevelChangeIncrease", "dexterity",
		2, observedLoggerSugared)
	assert.Equal(t, 17, a.Values["dexterity"])
	a.AdjustValues("AdditionalBonus", "strength",
		2, observedLoggerSugared)
	assert.Equal(t, 15, a.Values["strength"])
	actual, _ := a.GetModifier("strength")
	assert.Equal(t, 2, actual)
	actual, _ = a.GetModifier("intelligence")
	assert.Equal(t, -1, actual)
	fmt.Println(a.ToPrettyString())

	allLogs := observedLogs.All()
	assert.Equal(t, "AdjustValues", allLogs[len(allLogs)-1].Message)
}
