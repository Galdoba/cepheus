package interpolate

func MassByIndex(i int) float64 {
	massMap := make(map[int]float64)

	massMap[110] = 200
	massMap[115] = 80
	massMap[120] = 60
	massMap[125] = 30
	massMap[130] = 20
	massMap[135] = 15
	massMap[140] = 13
	massMap[145] = 12
	massMap[150] = 12
	massMap[155] = 13
	massMap[160] = 14
	massMap[165] = 18
	massMap[170] = 20
	massMap[175] = 25
	massMap[179] = 30

	massMap[210] = 150
	massMap[215] = 60
	massMap[220] = 40
	massMap[225] = 25
	massMap[230] = 15
	massMap[235] = 13
	massMap[240] = 12
	massMap[245] = 10
	massMap[250] = 10
	massMap[255] = 11
	massMap[260] = 12
	massMap[265] = 13
	massMap[270] = 15
	massMap[275] = 20
	massMap[279] = 25

	massMap[310] = 130
	massMap[315] = 40
	massMap[320] = 30
	massMap[325] = 20
	massMap[330] = 14
	massMap[335] = 11
	massMap[340] = 10
	massMap[345] = 8
	massMap[350] = 8
	massMap[355] = 10
	massMap[360] = 10
	massMap[365] = 12
	massMap[370] = 14
	massMap[375] = 16
	massMap[379] = 18

	massMap[410] = 110
	massMap[415] = 30
	massMap[420] = 20
	massMap[425] = 10
	massMap[430] = 8
	massMap[435] = 6
	massMap[440] = 4
	massMap[445] = 3
	massMap[450] = 2.5
	massMap[455] = 2.4
	massMap[460] = 1.1
	massMap[465] = 1.5
	massMap[470] = 1.8
	massMap[475] = 2.4
	massMap[479] = 8

	massMap[520] = 20
	massMap[525] = 10
	massMap[530] = 4
	massMap[535] = 2.3
	massMap[540] = 2
	massMap[545] = 1.5
	massMap[550] = 1.7
	massMap[555] = 1.2
	massMap[560] = 1.5

	massMap[610] = 90
	massMap[615] = 60
	massMap[620] = 18
	massMap[625] = 5
	massMap[630] = 2.2
	massMap[635] = 1.8
	massMap[640] = 1.5
	massMap[645] = 1.3
	massMap[650] = 1.1
	massMap[655] = 0.9
	massMap[660] = 0.8
	massMap[665] = 0.7
	massMap[670] = 0.5
	massMap[675] = 0.16
	massMap[679] = 0.08

	massMap[710] = 2
	massMap[715] = 1.5
	massMap[720] = 0.5
	massMap[725] = 0.4

	massMap[750] = 0.8
	massMap[755] = 0.7
	massMap[760] = 0.6
	massMap[765] = 0.5
	massMap[770] = 0.4
	massMap[775] = 0.12
	massMap[779] = 0.075

	massMap[910] = 0.08
	massMap[915] = 0.06
	massMap[920] = 0.05
	massMap[925] = 0.04
	massMap[930] = 0.025
	massMap[935] = 0.013

	return interpolate(massMap, i)
}

func TempByIndex(i int) float64 {
	tempMap := make(map[int]float64)
	if i < 900 {
		i = i % 100
	}
	tempMap[10] = 50000
	tempMap[15] = 40000
	tempMap[20] = 30000
	tempMap[25] = 15000
	tempMap[30] = 10000
	tempMap[35] = 8000
	tempMap[40] = 7500
	tempMap[45] = 6500
	tempMap[50] = 6000
	tempMap[55] = 5600
	tempMap[60] = 5200
	tempMap[65] = 4400
	tempMap[70] = 3700
	tempMap[75] = 3000
	tempMap[79] = 2400

	tempMap[910] = 2400
	tempMap[915] = 1850
	tempMap[920] = 1300
	tempMap[925] = 900
	tempMap[930] = 550
	tempMap[935] = 300
	return interpolate(tempMap, i)
}

func DiamByIndex(i int) float64 {
	diamMap := make(map[int]float64)

	diamMap[110] = 25
	diamMap[115] = 22
	diamMap[120] = 20
	diamMap[125] = 60
	diamMap[130] = 120
	diamMap[135] = 180
	diamMap[140] = 210
	diamMap[145] = 280
	diamMap[150] = 330
	diamMap[155] = 360
	diamMap[160] = 420
	diamMap[165] = 600
	diamMap[170] = 900
	diamMap[175] = 1200
	diamMap[179] = 1800

	diamMap[210] = 24
	diamMap[215] = 20
	diamMap[220] = 14
	diamMap[225] = 25
	diamMap[230] = 50
	diamMap[235] = 75
	diamMap[240] = 85
	diamMap[245] = 115
	diamMap[250] = 135
	diamMap[255] = 150
	diamMap[260] = 180
	diamMap[265] = 260
	diamMap[270] = 380
	diamMap[275] = 600
	diamMap[279] = 800

	diamMap[310] = 22
	diamMap[315] = 18
	diamMap[320] = 12
	diamMap[325] = 14
	diamMap[330] = 30
	diamMap[335] = 45
	diamMap[340] = 50
	diamMap[345] = 66
	diamMap[350] = 77
	diamMap[355] = 90
	diamMap[360] = 110
	diamMap[365] = 160
	diamMap[370] = 230
	diamMap[375] = 350
	diamMap[379] = 500

	diamMap[410] = 21
	diamMap[415] = 15
	diamMap[420] = 10
	diamMap[425] = 6
	diamMap[430] = 5
	diamMap[435] = 5
	diamMap[440] = 5
	diamMap[445] = 5
	diamMap[450] = 10
	diamMap[455] = 15
	diamMap[460] = 20
	diamMap[465] = 40
	diamMap[470] = 60
	diamMap[475] = 100
	diamMap[479] = 200

	// massMap[510] = 110
	// massMap[515] = 30
	diamMap[520] = 8
	diamMap[525] = 5
	diamMap[530] = 4
	diamMap[535] = 3
	diamMap[540] = 3
	diamMap[545] = 2
	diamMap[550] = 3
	diamMap[555] = 4
	diamMap[560] = 6
	// massMap[565] = 1.5
	// massMap[570] = 1.8
	// massMap[575] = 2.4
	// massMap[579] = 8

	diamMap[610] = 20
	diamMap[615] = 12
	diamMap[620] = 7
	diamMap[625] = 3.5
	diamMap[630] = 2.2
	diamMap[635] = 2
	diamMap[640] = 1.7
	diamMap[645] = 1.5
	diamMap[650] = 1.1
	diamMap[655] = 0.95
	diamMap[660] = 0.9
	diamMap[665] = 0.8
	diamMap[670] = 0.7
	diamMap[675] = 0.2
	diamMap[679] = 0.1

	diamMap[710] = 0.18
	diamMap[715] = 0.18
	diamMap[720] = 0.2
	diamMap[725] = 0.5
	// massMap[730] = 2.2
	// massMap[735] = 1.8
	// massMap[740] = 1.5
	// massMap[745] = 1.3
	diamMap[750] = 0.8
	diamMap[755] = 0.7
	diamMap[760] = 0.6
	diamMap[765] = 0.5
	diamMap[770] = 0.4
	diamMap[775] = 0.1
	diamMap[779] = 0.08

	diamMap[910] = 0.1
	diamMap[915] = 0.08
	diamMap[920] = 0.9
	diamMap[925] = 0.11
	diamMap[930] = 0.1
	diamMap[935] = 0.1
	return interpolate(diamMap, i)
}

func MAO_ByIndex(i int) float64 {
	maoMap := make(map[int]float64)

	maoMap[110] = 0.63
	maoMap[115] = 0.55
	maoMap[120] = 0.5
	maoMap[125] = 1.67
	maoMap[130] = 3.34
	maoMap[135] = 4.17
	maoMap[140] = 4.42
	maoMap[145] = 5.0
	maoMap[150] = 5.21
	maoMap[155] = 5.34
	maoMap[160] = 5.59
	maoMap[165] = 6.17
	maoMap[170] = 6.8
	maoMap[175] = 7.2
	maoMap[179] = 7.8

	maoMap[210] = 0.6
	maoMap[215] = 0.5
	maoMap[220] = 0.35
	maoMap[225] = 0.63
	maoMap[230] = 1.4
	maoMap[235] = 2.17
	maoMap[240] = 2.5
	maoMap[245] = 3.25
	maoMap[250] = 3.59
	maoMap[255] = 3.84
	maoMap[260] = 4.17
	maoMap[265] = 4.84
	maoMap[270] = 5.42
	maoMap[275] = 6.17
	maoMap[279] = 6.59

	maoMap[310] = 0.55
	maoMap[315] = 0.45
	maoMap[320] = 0.3
	maoMap[325] = 0.35
	maoMap[330] = 0.75
	maoMap[335] = 1.17
	maoMap[340] = 1.33
	maoMap[345] = 1.87
	maoMap[350] = 2.24
	maoMap[355] = 2.67
	maoMap[360] = 3.17
	maoMap[365] = 4.0
	maoMap[370] = 4.59
	maoMap[375] = 5.3
	maoMap[379] = 5.92

	maoMap[410] = 0.53
	maoMap[415] = 0.38
	maoMap[420] = 0.25
	maoMap[425] = 0.15
	maoMap[430] = 0.13
	maoMap[435] = 0.13
	maoMap[440] = 0.13
	maoMap[445] = 0.13
	maoMap[450] = 0.25
	maoMap[455] = 0.38
	maoMap[460] = 0.5
	maoMap[465] = 1.0
	maoMap[470] = 1.68
	maoMap[475] = 3.0
	maoMap[479] = 4.34

	// massMap[510] = 110
	// massMap[515] = 30
	maoMap[520] = 0.2
	maoMap[525] = 0.13
	maoMap[530] = 0.1
	maoMap[535] = 0.07
	maoMap[540] = 0.07
	maoMap[545] = 0.06
	maoMap[550] = 0.07
	maoMap[555] = 0.1
	maoMap[560] = 0.15
	// massMap[565] = 1.5
	// massMap[570] = 1.8
	// massMap[575] = 2.4
	// massMap[579] = 8

	maoMap[610] = 0.5
	maoMap[615] = 0.3
	maoMap[620] = 0.18
	maoMap[625] = 0.09
	maoMap[630] = 0.06
	maoMap[635] = 0.05
	maoMap[640] = 0.04
	maoMap[645] = 0.03
	maoMap[650] = 0.03
	maoMap[655] = 0.02
	maoMap[660] = 0.02
	maoMap[665] = 0.02
	maoMap[670] = 0.02
	maoMap[675] = 0.01
	maoMap[679] = 0.01

	maoMap[710] = 0.01
	maoMap[715] = 0.01
	maoMap[720] = 0.01
	maoMap[725] = 0.01
	// massMap[730] = 2.2
	// massMap[735] = 1.8
	// massMap[740] = 1.5
	// massMap[745] = 1.3
	maoMap[750] = 0.02
	maoMap[755] = 0.02
	maoMap[760] = 0.02
	maoMap[765] = 0.01
	maoMap[770] = 0.01
	maoMap[775] = 0.01
	maoMap[779] = 0.01

	maoMap[910] = 0.005
	maoMap[915] = 0.005
	maoMap[920] = 0.005
	maoMap[925] = 0.005
	maoMap[930] = 0.005
	maoMap[935] = 0.005
	return interpolate(maoMap, i)
}
