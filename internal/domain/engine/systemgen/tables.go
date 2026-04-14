package systemgen

// ---------------------------------------------------------------------------
// Step 1: Subsector object presence thresholds
// ---------------------------------------------------------------------------

// objectPresenceThreshold returns the d100 value at or below which an object
// exists in a hex for the given subsector density.
var objectPresenceThreshold = map[SubsectorType]int{
	SubEmpty:     5,
	SubScattered: 20,
	SubDispersed: 35,
	SubAverage:   50,
	SubCrowded:   60,
	SubDense:     75,
}

// ---------------------------------------------------------------------------
// Step 2: Object type determination (d100)
// ---------------------------------------------------------------------------

type objectRange struct {
	Min int
	Max int
	Type ObjectType
}

var objectTypeRanges = []objectRange{
	{Min: 1, Max: 80, Type: ObjectStar},
	{Min: 81, Max: 88, Type: ObjectBrownDwarf},
	{Min: 89, Max: 94, Type: ObjectRoguePlanet},
	{Min: 95, Max: 97, Type: ObjectRogueGasGiant},
	{Min: 98, Max: 98, Type: ObjectNeutronStar},
	{Min: 99, Max: 99, Type: ObjectNebula},
	{Min: 100, Max: 100, Type: ObjectBlackHole}, // d100 returns 1-100; "00" maps to 100
}

// ---------------------------------------------------------------------------
// Step 3: Brown dwarf class (d100)
// ---------------------------------------------------------------------------

var brownDwarfClassRanges = []struct {
	Min   int
	Max   int
	Class BrownDwarfClass
}{
	{Min: 1, Max: 50, Class: ClassL},
	{Min: 51, Max: 75, Class: ClassT},
	{Min: 76, Max: 100, Class: ClassY},
}

// ---------------------------------------------------------------------------
// Step 4: Stellar class tables (d100)
// ---------------------------------------------------------------------------

type stellarRange struct {
	Min   int
	Max   int
	Class StellarClass
}

var stellarRealistic = []stellarRange{
	{Min: 1, Max: 80, Class: ClassM},
	{Min: 81, Max: 88, Class: ClassK},
	{Min: 89, Max: 94, Class: ClassG},
	{Min: 95, Max: 97, Class: ClassF},
	{Min: 98, Max: 98, Class: ClassA},
	{Min: 99, Max: 99, Class: ClassB},
	{Min: 100, Max: 100, Class: ClassO},
}

var stellarSemiRealistic = []stellarRange{
	{Min: 1, Max: 50, Class: ClassM},
	{Min: 51, Max: 77, Class: ClassK},
	{Min: 78, Max: 90, Class: ClassG},
	{Min: 91, Max: 97, Class: ClassF},
	{Min: 98, Max: 98, Class: ClassA},
	{Min: 99, Max: 99, Class: ClassB},
	{Min: 100, Max: 100, Class: ClassO},
}

var stellarFantastic = []stellarRange{
	{Min: 1, Max: 25, Class: ClassM},
	{Min: 26, Max: 50, Class: ClassK},
	{Min: 51, Max: 75, Class: ClassG},
	{Min: 76, Max: 97, Class: ClassF},
	{Min: 98, Max: 98, Class: ClassA},
	{Min: 99, Max: 99, Class: ClassB},
	{Min: 100, Max: 100, Class: ClassO},
}

// ---------------------------------------------------------------------------
// Step 6: Luminosity class (d100)
// ---------------------------------------------------------------------------

type luminosityRange struct {
	Min   int
	Max   int
	Class LuminosityClass
}

var luminosityClassRanges = []luminosityRange{
	{Min: 1, Max: 90, Class: LumV},
	{Min: 91, Max: 94, Class: LumIV},
	{Min: 95, Max: 96, Class: LumD},
	{Min: 97, Max: 99, Class: LumIII},
	// 00 → sub-roll on d10 table
}

// Sub-roll for luminosity 00 (d10, where 0 maps to 10):
var luminositySubRollRanges = []luminosityRange{
	{Min: 1, Max: 4, Class: LumII},
	{Min: 5, Max: 5, Class: LumVI},
	{Min: 7, Max: 8, Class: LumIa},
	{Min: 9, Max: 10, Class: LumIb},
}

// ---------------------------------------------------------------------------
// Step 7: Multiple star system tables (d100)
// ---------------------------------------------------------------------------

var multipleStarOBO = []struct {
	Min  int
	Max  int
	Count int
}{
	{Min: 1, Max: 10, Count: 1},
	{Min: 11, Max: 90, Count: 2},
	{Min: 91, Max: 98, Count: 3},
	{Min: 99, Max: 99, Count: 4},
	{Min: 100, Max: 100, Count: 5},
}

var multipleStarKGF = []struct {
	Min  int
	Max  int
	Count int
}{
	{Min: 1, Max: 45, Count: 1},
	{Min: 46, Max: 99, Count: 2},
	{Min: 100, Max: 100, Count: 3},
}

var multipleStarM = []struct {
	Min  int
	Max  int
	Count int
}{
	{Min: 1, Max: 69, Count: 1},
	{Min: 70, Max: 98, Count: 2},
	{Min: 99, Max: 100, Count: 3},
}

// ---------------------------------------------------------------------------
// Step 8: Companion distance tables
// ---------------------------------------------------------------------------

type companionDistRange struct {
	Min      int
	Max      int
	Distance CompanionDistance
}

var companionDistanceRanges = []companionDistRange{
	{Min: 1, Max: 10, Distance: DistContact},
	{Min: 11, Max: 30, Distance: DistClose},
	{Min: 31, Max: 50, Distance: DistNear},
	{Min: 51, Max: 80, Distance: DistFar},
	{Min: 81, Max: 100, Distance: DistDistant},
}

// AU sub-tables: each entry is {d100 range min, d100 range max, AU value}
type auRange struct {
	Min int
	Max int
	AU  float64
}

var closeCompanionAU = []auRange{
	{Min: 1, Max: 9, AU: 0.5},
	{Min: 10, Max: 19, AU: 1.5},
	{Min: 20, Max: 29, AU: 2.0},
	{Min: 30, Max: 39, AU: 2.5},
	{Min: 40, Max: 49, AU: 3.0},
	{Min: 50, Max: 59, AU: 3.5},
	{Min: 60, Max: 69, AU: 4.0},
	{Min: 70, Max: 79, AU: 4.5},
	{Min: 80, Max: 89, AU: 5.0},
	{Min: 90, Max: 100, AU: 5.5},
}

var nearCompanionAU = []auRange{
	{Min: 1, Max: 9, AU: 10},
	{Min: 10, Max: 19, AU: 20},
	{Min: 20, Max: 29, AU: 30},
	{Min: 30, Max: 39, AU: 40},
	{Min: 40, Max: 49, AU: 50},
	{Min: 50, Max: 59, AU: 60},
	{Min: 60, Max: 69, AU: 70},
	{Min: 70, Max: 79, AU: 80},
	{Min: 80, Max: 89, AU: 90},
	{Min: 90, Max: 100, AU: 100},
}

var farCompanionAU = []auRange{
	{Min: 1, Max: 9, AU: 100},
	{Min: 10, Max: 19, AU: 150},
	{Min: 20, Max: 29, AU: 200},
	{Min: 30, Max: 39, AU: 250},
	{Min: 40, Max: 49, AU: 300},
	{Min: 50, Max: 59, AU: 350},
	{Min: 60, Max: 69, AU: 400},
	{Min: 70, Max: 79, AU: 450},
	{Min: 80, Max: 89, AU: 500},
	{Min: 90, Max: 100, AU: 550},
}

var distantCompanionAU = []auRange{
	{Min: 1, Max: 9, AU: 600},
	{Min: 10, Max: 19, AU: 750},
	{Min: 20, Max: 29, AU: 1000},
	{Min: 30, Max: 39, AU: 1500},
	{Min: 40, Max: 49, AU: 2000},
	{Min: 50, Max: 59, AU: 2500},
	{Min: 60, Max: 69, AU: 3000},
	{Min: 70, Max: 79, AU: 4000},
	{Min: 80, Max: 89, AU: 5000},
	{Min: 90, Max: 100, AU: 6000},
}

// ---------------------------------------------------------------------------
// Step 11: Star zone data tables
//
// Format: map[numericClass]StarZoneEntry
// Key: 0-9 for numeric classification
// ---------------------------------------------------------------------------

// Main Sequence (V)
var mainSequenceZones = map[int]StarZoneEntry{
	0: {Name: "O0 V", Temperature: 50000, Mass: 100, Luminosity: 1240000, InnerLimit: 20, HabitableMin: 1057.88, HabitableMax: 1447.62, SnowLine: 0, OuterLimit: 4000},
	1: {Name: "O1 V", Temperature: 47800, Mass: 97.5, Luminosity: 994000, InnerLimit: 19.5, HabitableMin: 947.15, HabitableMax: 1296.09, SnowLine: 0, OuterLimit: 3900},
	2: {Name: "O2 V", Temperature: 45600, Mass: 95, Luminosity: 795000, InnerLimit: 19, HabitableMin: 847.05, HabitableMax: 1159.12, SnowLine: 0, OuterLimit: 3800},
	3: {Name: "O3 V", Temperature: 43400, Mass: 92.5, Luminosity: 634000, InnerLimit: 18.5, HabitableMin: 756.43, HabitableMax: 1035.11, SnowLine: 0, OuterLimit: 3700},
	4: {Name: "O4 V", Temperature: 41200, Mass: 90, Luminosity: 504000, InnerLimit: 18, HabitableMin: 674.43, HabitableMax: 922.91, SnowLine: 3549.65, OuterLimit: 3600},
	5: {Name: "O5 V", Temperature: 39000, Mass: 60, Luminosity: 398000, InnerLimit: 12, HabitableMin: 599.33, HabitableMax: 820.13, SnowLine: 0, OuterLimit: 2400},
	6: {Name: "O6 V", Temperature: 36800, Mass: 37, Luminosity: 260000, InnerLimit: 7.4, HabitableMin: 153.19, HabitableMax: 662.87, SnowLine: 0, OuterLimit: 1480},
	7: {Name: "O7 V", Temperature: 34600, Mass: 30, Luminosity: 154000, InnerLimit: 6, HabitableMin: 372.81, HabitableMax: 510.16, SnowLine: 0, OuterLimit: 1200},
	8: {Name: "O8 V", Temperature: 32400, Mass: 23, Luminosity: 99100, InnerLimit: 4.6, HabitableMin: 299.06, HabitableMax: 409.24, SnowLine: 0, OuterLimit: 920},
	9: {Name: "O9 V", Temperature: 30200, Mass: 20, Luminosity: 57600, InnerLimit: 4, HabitableMin: 228, HabitableMax: 312, SnowLine: 0, OuterLimit: 800},
}

// B main sequence (separate map to avoid key collision — keyed differently)
var mainSequenceB = map[int]StarZoneEntry{
	0: {Name: "B0 V", Temperature: 28000, Mass: 17.5, Luminosity: 36200, InnerLimit: 3.5, HabitableMin: 180.75, HabitableMax: 247.34, SnowLine: 0, OuterLimit: 700},
	1: {Name: "B1 V", Temperature: 26190, Mass: 14.2, Luminosity: 19400, InnerLimit: 2.84, HabitableMin: 132.32, HabitableMax: 181.07, SnowLine: 0, OuterLimit: 568},
	2: {Name: "B2 V", Temperature: 24380, Mass: 10.9, Luminosity: 9360, InnerLimit: 2.18, HabitableMin: 89.97, HabitableMax: 125.77, SnowLine: 0, OuterLimit: 436},
	3: {Name: "B3 V", Temperature: 22570, Mass: 7.6, Luminosity: 4890, InnerLimit: 1.52, HabitableMin: 66.43, HabitableMax: 90.91, SnowLine: 0, OuterLimit: 304},
	4: {Name: "B4 V", Temperature: 20760, Mass: 6.7, Luminosity: 2290, InnerLimit: 1.34, HabitableMin: 45.46, HabitableMax: 62.21, SnowLine: 239.27, OuterLimit: 268},
	5: {Name: "B5 V", Temperature: 18950, Mass: 5.9, Luminosity: 1160, InnerLimit: 1.18, HabitableMin: 32.36, HabitableMax: 44.28, SnowLine: 170.29, OuterLimit: 236},
	6: {Name: "B6 V", Temperature: 17140, Mass: 5.2, Luminosity: 692, InnerLimit: 1.04, HabitableMin: 24.99, HabitableMax: 34.20, SnowLine: 131.53, OuterLimit: 208},
	7: {Name: "B7 V", Temperature: 15330, Mass: 4.5, Luminosity: 404, InnerLimit: 0.90, HabitableMin: 19.09, HabitableMax: 26.13, SnowLine: 100.50, OuterLimit: 180},
	8: {Name: "B8 V", Temperature: 13520, Mass: 3.8, Luminosity: 211, InnerLimit: 0.76, HabitableMin: 13.80, HabitableMax: 18.88, SnowLine: 76.63, OuterLimit: 152},
	9: {Name: "B9 V", Temperature: 11710, Mass: 3.4, Luminosity: 119, InnerLimit: 0.68, HabitableMin: 10.36, HabitableMax: 14.18, SnowLine: 54.54, OuterLimit: 136},
}

var mainSequenceA = map[int]StarZoneEntry{
	0: {Name: "A0 V", Temperature: 9900, Mass: 2.9, Luminosity: 67.4, InnerLimit: 0.58, HabitableMin: 7.80, HabitableMax: 10.67, SnowLine: 41.05, OuterLimit: 116},
	1: {Name: "A1 V", Temperature: 9650, Mass: 2.7, Luminosity: 49.2, InnerLimit: 0.54, HabitableMin: 6.66, HabitableMax: 9.12, SnowLine: 35.07, OuterLimit: 108},
	2: {Name: "A2 V", Temperature: 9400, Mass: 2.5, Luminosity: 39.4, InnerLimit: 0.50, HabitableMin: 5.96, HabitableMax: 8.16, SnowLine: 31.38, OuterLimit: 100},
	3: {Name: "A3 V", Temperature: 9150, Mass: 2.4, Luminosity: 28.9, InnerLimit: 0.48, HabitableMin: 5.11, HabitableMax: 6.99, SnowLine: 26.88, OuterLimit: 96},
	4: {Name: "A4 V", Temperature: 8900, Mass: 2.1, Luminosity: 23.2, InnerLimit: 0.42, HabitableMin: 4.58, HabitableMax: 6.26, SnowLine: 24.08, OuterLimit: 84},
	5: {Name: "A5 V", Temperature: 8650, Mass: 1.9, Luminosity: 17.0, InnerLimit: 0.38, HabitableMin: 3.92, HabitableMax: 5.36, SnowLine: 20.62, OuterLimit: 76},
	6: {Name: "A6 V", Temperature: 8400, Mass: 1.8, Luminosity: 15.1, InnerLimit: 0.36, HabitableMin: 3.69, HabitableMax: 5.05, SnowLine: 19.43, OuterLimit: 72},
	7: {Name: "A7 V", Temperature: 8150, Mass: 1.8, Luminosity: 12.2, InnerLimit: 0.36, HabitableMin: 3.32, HabitableMax: 4.54, SnowLine: 17.46, OuterLimit: 72},
	8: {Name: "A8 V", Temperature: 7900, Mass: 1.8, Luminosity: 10.9, InnerLimit: 0.36, HabitableMin: 3.14, HabitableMax: 4.30, SnowLine: 16.50, OuterLimit: 72},
	9: {Name: "A9 V", Temperature: 7650, Mass: 1.7, Luminosity: 8.85, InnerLimit: 0.34, HabitableMin: 2.83, HabitableMax: 3.87, SnowLine: 14.87, OuterLimit: 68},
}

var mainSequenceF = map[int]StarZoneEntry{
	0: {Name: "F0 V", Temperature: 7400, Mass: 1.6, Luminosity: 7.94, InnerLimit: 0.32, HabitableMin: 2.68, HabitableMax: 3.66, SnowLine: 14.09, OuterLimit: 64},
	1: {Name: "F1 V", Temperature: 7260, Mass: 1.6, Luminosity: 6.56, InnerLimit: 0.32, HabitableMin: 2.43, HabitableMax: 3.17, SnowLine: 12.81, OuterLimit: 64},
	2: {Name: "F2 V", Temperature: 7120, Mass: 1.5, Luminosity: 5.95, InnerLimit: 0.30, HabitableMin: 2.32, HabitableMax: 3.17, SnowLine: 12.20, OuterLimit: 60},
	3: {Name: "F3 V", Temperature: 6980, Mass: 1.5, Luminosity: 4.94, InnerLimit: 0.30, HabitableMin: 2.11, HabitableMax: 2.89, SnowLine: 11.11, OuterLimit: 60},
	4: {Name: "F4 V", Temperature: 6840, Mass: 1.4, Luminosity: 4.50, InnerLimit: 0.28, HabitableMin: 2.02, HabitableMax: 2.76, SnowLine: 10.61, OuterLimit: 56},
	5: {Name: "F5 V", Temperature: 6700, Mass: 1.4, Luminosity: 3.75, InnerLimit: 0.28, HabitableMin: 1.84, HabitableMax: 2.52, SnowLine: 9.68, OuterLimit: 56},
	6: {Name: "F6 V", Temperature: 6560, Mass: 1.3, Luminosity: 3.13, InnerLimit: 0.26, HabitableMin: 1.68, HabitableMax: 2.30, SnowLine: 8.85, OuterLimit: 52},
	7: {Name: "F7 V", Temperature: 6420, Mass: 1.3, Luminosity: 2.62, InnerLimit: 0.26, HabitableMin: 1.54, HabitableMax: 2.10, SnowLine: 8.09, OuterLimit: 52},
	8: {Name: "F8 V", Temperature: 6280, Mass: 1.2, Luminosity: 2.41, InnerLimit: 0.24, HabitableMin: 1.47, HabitableMax: 2.02, SnowLine: 7.76, OuterLimit: 48},
	9: {Name: "F9 V", Temperature: 6140, Mass: 1.1, Luminosity: 2.03, InnerLimit: 0.22, HabitableMin: 1.35, HabitableMax: 1.85, SnowLine: 7.12, OuterLimit: 44},
}

var mainSequenceG = map[int]StarZoneEntry{
	0: {Name: "G0 V", Temperature: 6000, Mass: 1.1, Luminosity: 1.72, InnerLimit: 0.22, HabitableMin: 1.25, HabitableMax: 1.70, SnowLine: 6.56, OuterLimit: 44},
	1: {Name: "G1 V", Temperature: 5890, Mass: 1.0, Luminosity: 1.46, InnerLimit: 0.20, HabitableMin: 1.15, HabitableMax: 1.57, SnowLine: 6.04, OuterLimit: 40},
	2: {Name: "G2 V", Temperature: 5780, Mass: 1.0, Luminosity: 1.00, InnerLimit: 0.20, HabitableMin: 0.95, HabitableMax: 1.30, SnowLine: 5.00, OuterLimit: 40},
	3: {Name: "G3 V", Temperature: 5670, Mass: 1.0, Luminosity: 1.00, InnerLimit: 0.20, HabitableMin: 0.95, HabitableMax: 1.30, SnowLine: 5.00, OuterLimit: 40},
	4: {Name: "G4 V", Temperature: 5560, Mass: 0.9, Luminosity: 0.98, InnerLimit: 0.18, HabitableMin: 0.94, HabitableMax: 1.29, SnowLine: 4.95, OuterLimit: 36},
	5: {Name: "G5 V", Temperature: 5450, Mass: 0.9, Luminosity: 0.84, InnerLimit: 0.18, HabitableMin: 0.87, HabitableMax: 1.19, SnowLine: 4.58, OuterLimit: 36},
	6: {Name: "G6 V", Temperature: 5340, Mass: 0.9, Luminosity: 0.79, InnerLimit: 0.18, HabitableMin: 0.84, HabitableMax: 1.16, SnowLine: 4.44, OuterLimit: 36},
	7: {Name: "G7 V", Temperature: 5230, Mass: 0.9, Luminosity: 0.68, InnerLimit: 0.18, HabitableMin: 0.78, HabitableMax: 1.04, SnowLine: 4.24, OuterLimit: 35},
}

var mainSequenceK = map[int]StarZoneEntry{
	0: {Name: "K0 V", Temperature: 4900, Mass: 0.8, Luminosity: 0.54, InnerLimit: 0.16, HabitableMin: 0.70, HabitableMax: 0.96, SnowLine: 3.67, OuterLimit: 32},
	1: {Name: "K1 V", Temperature: 4760, Mass: 0.8, Luminosity: 0.44, InnerLimit: 0.16, HabitableMin: 0.63, HabitableMax: 0.86, SnowLine: 3.32, OuterLimit: 32},
	2: {Name: "K2 V", Temperature: 4620, Mass: 0.7, Luminosity: 0.40, InnerLimit: 0.14, HabitableMin: 0.60, HabitableMax: 0.82, SnowLine: 3.16, OuterLimit: 28},
	3: {Name: "K3 V", Temperature: 4480, Mass: 0.7, Luminosity: 0.34, InnerLimit: 0.14, HabitableMin: 0.55, HabitableMax: 0.76, SnowLine: 2.92, OuterLimit: 28},
	4: {Name: "K4 V", Temperature: 4340, Mass: 0.7, Luminosity: 0.31, InnerLimit: 0.14, HabitableMin: 0.53, HabitableMax: 0.72, SnowLine: 2.78, OuterLimit: 28},
	5: {Name: "K5 V", Temperature: 4200, Mass: 0.7, Luminosity: 0.27, InnerLimit: 0.14, HabitableMin: 0.49, HabitableMax: 0.68, SnowLine: 2.60, OuterLimit: 28},
	6: {Name: "K6 V", Temperature: 4060, Mass: 0.6, Luminosity: 0.21, InnerLimit: 0.12, HabitableMin: 0.44, HabitableMax: 0.60, SnowLine: 2.29, OuterLimit: 24},
	7: {Name: "K7 V", Temperature: 3920, Mass: 0.6, Luminosity: 0.19, InnerLimit: 0.12, HabitableMin: 0.41, HabitableMax: 0.57, SnowLine: 2.18, OuterLimit: 24},
	8: {Name: "K8 V", Temperature: 3780, Mass: 0.6, Luminosity: 0.16, InnerLimit: 0.12, HabitableMin: 0.38, HabitableMax: 0.52, SnowLine: 2.00, OuterLimit: 24},
	9: {Name: "K9 V", Temperature: 3640, Mass: 0.5, Luminosity: 0.14, InnerLimit: 0.10, HabitableMin: 0.36, HabitableMax: 0.49, SnowLine: 1.87, OuterLimit: 20},
}

var mainSequenceM = map[int]StarZoneEntry{
	0: {Name: "M0 V", Temperature: 3500, Mass: 0.5, Luminosity: 0.125, InnerLimit: 0.100, HabitableMin: 0.336, HabitableMax: 0.460, SnowLine: 1.768, OuterLimit: 20},
	1: {Name: "M1 V", Temperature: 3333, Mass: 0.5, Luminosity: 0.0618, InnerLimit: 0.100, HabitableMin: 0.236, HabitableMax: 0.323, SnowLine: 1.243, OuterLimit: 20},
	2: {Name: "M2 V", Temperature: 3167, Mass: 0.4, Luminosity: 0.0321, InnerLimit: 0.080, HabitableMin: 0.170, HabitableMax: 0.233, SnowLine: 0.896, OuterLimit: 16},
	3: {Name: "M3 V", Temperature: 3000, Mass: 0.3, Luminosity: 0.0178, InnerLimit: 0.060, HabitableMin: 0.127, HabitableMax: 0.173, SnowLine: 0.667, OuterLimit: 12},
	4: {Name: "M4 V", Temperature: 2833, Mass: 0.3, Luminosity: 0.0106, InnerLimit: 0.060, HabitableMin: 0.098, HabitableMax: 0.134, SnowLine: 0.515, OuterLimit: 12},
	5: {Name: "M5 V", Temperature: 2667, Mass: 0.2, Luminosity: 0.00624, InnerLimit: 0.040, HabitableMin: 0.075, HabitableMax: 0.103, SnowLine: 0.395, OuterLimit: 8},
	6: {Name: "M6 V", Temperature: 2500, Mass: 0.2, Luminosity: 0.00450, InnerLimit: 0.040, HabitableMin: 0.0637, HabitableMax: 0.0872, SnowLine: 0.335, OuterLimit: 8},
	7: {Name: "M7 V", Temperature: 2333, Mass: 0.1, Luminosity: 0.00369, InnerLimit: 0.020, HabitableMin: 0.0577, HabitableMax: 0.0790, SnowLine: 0.960, OuterLimit: 4},
	8: {Name: "M8 V", Temperature: 2167, Mass: 0.1, Luminosity: 0.00353, InnerLimit: 0.020, HabitableMin: 0.0564, HabitableMax: 0.0772, SnowLine: 0.297, OuterLimit: 4},
	9: {Name: "M9 V", Temperature: 2000, Mass: 0.1, Luminosity: 0.00315, InnerLimit: 0.020, HabitableMin: 0.0533, HabitableMax: 0.0730, SnowLine: 0.281, OuterLimit: 4},
}

// Brown Dwarfs
var brownDwarfZones = map[string]StarZoneEntry{
	"L0": {Name: "L0", Temperature: 2200, Mass: 0.08, Luminosity: 0.005, InnerLimit: 0.016, HabitableMin: 0.067, HabitableMax: 0.092, SnowLine: 0.354, OuterLimit: 3.2},
	"L1": {Name: "L1", Temperature: 2100, Mass: 0.08, Luminosity: 0.004, InnerLimit: 0.016, HabitableMin: 0.060, HabitableMax: 0.082, SnowLine: 0.316, OuterLimit: 3.2},
	"L2": {Name: "L2", Temperature: 2000, Mass: 0.07, Luminosity: 0.003, InnerLimit: 0.014, HabitableMin: 0.052, HabitableMax: 0.071, SnowLine: 0.274, OuterLimit: 2.8},
	"L3": {Name: "L3", Temperature: 1900, Mass: 0.07, Luminosity: 0.001, InnerLimit: 0.014, HabitableMin: 0.030, HabitableMax: 0.041, SnowLine: 0.158, OuterLimit: 2.8},
	"L4": {Name: "L4", Temperature: 1800, Mass: 0.07, Luminosity: 0.0007, InnerLimit: 0.014, HabitableMin: 0.025, HabitableMax: 0.034, SnowLine: 0.132, OuterLimit: 2.8},
	"L5": {Name: "L5", Temperature: 1700, Mass: 0.06, Luminosity: 0.0005, InnerLimit: 0.012, HabitableMin: 0.021, HabitableMax: 0.029, SnowLine: 0.112, OuterLimit: 2.4},
	"L6": {Name: "L6", Temperature: 1600, Mass: 0.06, Luminosity: 0.0001, InnerLimit: 0.012, HabitableMin: 0.012, HabitableMax: 0.013, SnowLine: 0.050, OuterLimit: 2.4},
	"L7": {Name: "L7", Temperature: 1450, Mass: 0.05, Luminosity: 0.00007, InnerLimit: 0.010, HabitableMin: 0.010, HabitableMax: 0.011, SnowLine: 0.418, OuterLimit: 2.0},
	"L8": {Name: "L8", Temperature: 1425, Mass: 0.05, Luminosity: 0.00005, InnerLimit: 0.010, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.035, OuterLimit: 2.0},
	"L9": {Name: "L9", Temperature: 1410, Mass: 0.05, Luminosity: 0.00001, InnerLimit: 0.010, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.016, OuterLimit: 2.0},
	"T0": {Name: "T0", Temperature: 1400, Mass: 0.05, Luminosity: 0.0000060, InnerLimit: 0.010, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.012, OuterLimit: 2.0},
	"T1": {Name: "T1", Temperature: 1350, Mass: 0.04, Luminosity: 0.0000060, InnerLimit: 0.008, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.012, OuterLimit: 1.6},
	"T2": {Name: "T2", Temperature: 1300, Mass: 0.04, Luminosity: 0.0000055, InnerLimit: 0.008, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.012, OuterLimit: 1.6},
	"T3": {Name: "T3", Temperature: 1200, Mass: 0.04, Luminosity: 0.0000050, InnerLimit: 0.008, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.011, OuterLimit: 1.6},
	"T4": {Name: "T4", Temperature: 1100, Mass: 0.04, Luminosity: 0.0000040, InnerLimit: 0.008, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.010, OuterLimit: 1.6},
	"T5": {Name: "T5", Temperature: 1000, Mass: 0.03, Luminosity: 0.0000040, InnerLimit: 0.006, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.010, OuterLimit: 1.2},
	"T6": {Name: "T6", Temperature: 900, Mass: 0.03, Luminosity: 0.0000035, InnerLimit: 0.006, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.009, OuterLimit: 1.2},
	"T7": {Name: "T7", Temperature: 800, Mass: 0.03, Luminosity: 0.0000030, InnerLimit: 0.006, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.009, OuterLimit: 1.2},
	"T8": {Name: "T8", Temperature: 750, Mass: 0.03, Luminosity: 0.0000020, InnerLimit: 0.006, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.007, OuterLimit: 1.2},
	"T9": {Name: "T9", Temperature: 700, Mass: 0.03, Luminosity: 0.0000010, InnerLimit: 0.006, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.006, OuterLimit: 1.2},
	"Y0": {Name: "Y0", Temperature: 448, Mass: 0.03, Luminosity: 0.0000006, InnerLimit: 0.006, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.006, OuterLimit: 1.2},
	"Y1": {Name: "Y1", Temperature: 433, Mass: 0.03, Luminosity: 0.0000006, InnerLimit: 0.006, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.006, OuterLimit: 1.2},
	"Y2": {Name: "Y2", Temperature: 418, Mass: 0.03, Luminosity: 0.0000005, InnerLimit: 0.006, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.006, OuterLimit: 1.2},
	"Y3": {Name: "Y3", Temperature: 403, Mass: 0.03, Luminosity: 0.0000003, InnerLimit: 0.006, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.006, OuterLimit: 1.2},
	"Y4": {Name: "Y4", Temperature: 388, Mass: 0.03, Luminosity: 0.0000001, InnerLimit: 0.006, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.006, OuterLimit: 1.2},
	"Y5": {Name: "Y5", Temperature: 373, Mass: 0.03, Luminosity: 0.00000007, InnerLimit: 0.006, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.006, OuterLimit: 1.2},
	"Y6": {Name: "Y6", Temperature: 358, Mass: 0.03, Luminosity: 0.00000005, InnerLimit: 0.006, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.006, OuterLimit: 1.2},
	"Y7": {Name: "Y7", Temperature: 343, Mass: 0.03, Luminosity: 0.00000003, InnerLimit: 0.006, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.006, OuterLimit: 1.2},
	"Y8": {Name: "Y8", Temperature: 328, Mass: 0.03, Luminosity: 0.00000001, InnerLimit: 0.006, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.006, OuterLimit: 1.2},
	"Y9": {Name: "Y9", Temperature: 298, Mass: 0.03, Luminosity: 0.00000007, InnerLimit: 0.006, HabitableMin: 0, HabitableMax: 0, SnowLine: 0.006, OuterLimit: 1.2},
}

// ---------------------------------------------------------------------------
// Step 13: Gas giant size tables
// ---------------------------------------------------------------------------

// Neptunian (small gas giant): 2d6 → size code
var neptunianSizes = map[int]int{
	2: 30, 3: 32, 4: 35, 5: 37, 6: 40,
	7: 42, 8: 45, 9: 47, 10: 50, 11: 55, 12: 57,
}

// Jovian (large gas giant): 2d10 → size code
var jovianSizes = map[int]int{
	2: 60, 3: 70, 4: 80, 5: 90, 6: 100,
	7: 110, 8: 120, 9: 130, 10: 140, 11: 150,
	12: 160, 13: 170, 14: 180, 15: 190, 16: 200,
	17: 210, 18: 220, 19: 230, 20: 240,
}

// ---------------------------------------------------------------------------
// Step 18: System quirk table (d66 → description)
//
// d66 = first d6 * 10 + second d6 (results 11-66)
// ---------------------------------------------------------------------------

type systemQuirkEntry struct {
	Roll        int
	Description string
}

var systemQuirkTable = []systemQuirkEntry{
	{Roll: 11, Description: "The primary star in this system is prone to powerful solar flares. Subtract 2 from the population code of every planet and moon in the system."},
	{Roll: 12, Description: "One of the asteroid belts in this system is much more densely packed than the average asteroid belt."},
	{Roll: 13, Description: "One of the rocky planets has a water core instead of rocky/molten core."},
	{Roll: 14, Description: "The primary star has a dense debris cloud which is 1.5 AU wide."},
	{Roll: 15, Description: "The system has an extremely dense Oort cloud. Comet activity is high."},
	{Roll: 16, Description: "A gas giant has migrated to the center of the habitable zone."},
	{Roll: 21, Description: "A rocky planet with atmosphere < 3 has polar ice (hydrographic code 1)."},
	{Roll: 22, Description: "The star is prone to large-scale flares. Subtract 6 from biology chart rolls."},
	{Roll: 23, Description: "An asteroid belt beyond the habitable zone is rich in metals."},
	{Roll: 24, Description: "The system has an additional 1d6 rocky planets."},
	{Roll: 25, Description: "A rocky planet beyond the snow line is covered in water under thick ice."},
	{Roll: 26, Description: "All outgoing Zimm Points are closer to the star than incoming Zimm Points."},
	{Roll: 31, Description: "A planet has an elliptical orbit crossing another world's path."},
	{Roll: 32, Description: "All Zimm Points are in the outer system beyond the Snow Line."},
	{Roll: 33, Description: "An asteroid belt contains a Dwarf or Mercurian planet within it."},
	{Roll: 34, Description: "A gas giant has a massive rotating storm."},
	{Roll: 35, Description: "A planet in or near the habitable zone has an extremely eccentric orbit."},
	{Roll: 36, Description: "The system has an additional 1d6-3 (min 1) gas giants."},
	{Roll: 41, Description: "A rocky planet has a moon one size larger than the maximum."},
	{Roll: 42, Description: "A rocky planet has a retrograde orbit around the star."},
	{Roll: 43, Description: "Life has taken hold unusually well. Add +5 to the biology chart."},
	{Roll: 44, Description: "There is an additional star in this system."},
	{Roll: 45, Description: "The system has a brown dwarf as an additional distant companion."},
	{Roll: 46, Description: "The system has an additional 1d6-3 (min 1) asteroid belts."},
	{Roll: 51, Description: "A rocky planet has a twin world on the opposite side of the star."},
	{Roll: 52, Description: "A planet has strayed too close to the star, siphoning its surface."},
	{Roll: 53, Description: "All worlds in the system have circular orbits (eccentricity 0)."},
	{Roll: 54, Description: "A world has an extreme orbit at an odd angle to the elliptical plane."},
	{Roll: 55, Description: "A planet between the inner limit and snow line has extreme axial tilt."},
	{Roll: 56, Description: "The main world has a ring system with 1d10 rings of 1d10 km width."},
	{Roll: 61, Description: "A planet has a substantially darker surface. Add 0.20 to albedo."},
	{Roll: 62, Description: "A rocky planet has a companion planet (binary/double planet)."},
	{Roll: 63, Description: "A gas giant has native life in its upper atmosphere."},
	{Roll: 64, Description: "The system has no rocky planets other than the main world."},
	{Roll: 65, Description: "A planet bears evidence of an ancient alien civilization."},
	{Roll: 66, Description: "There is an abandoned alien megastructure in the system."},
}
