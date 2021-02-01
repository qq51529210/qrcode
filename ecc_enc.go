package qrcode

var (
	// galois对数表
	galoisLogTable = []byte{
		0, 0, 1, 25, 2, 50, 26, 198, 3, 223,
		51, 238, 27, 104, 199, 75, 4, 100, 224, 14,
		52, 141, 239, 129, 28, 193, 105, 248, 200, 8,
		76, 113, 5, 138, 101, 47, 225, 36, 15, 33,
		53, 147, 142, 218, 240, 18, 130, 69, 29, 181,
		194, 125, 106, 39, 249, 185, 201, 154, 9, 120,
		77, 228, 114, 166, 6, 191, 139, 98, 102, 221,
		48, 253, 226, 152, 37, 179, 16, 145, 34, 136,
		54, 208, 148, 206, 143, 150, 219, 189, 241, 210,
		19, 92, 131, 56, 70, 64, 30, 66, 182, 163,
		195, 72, 126, 110, 107, 58, 40, 84, 250, 133,
		186, 61, 202, 94, 155, 159, 10, 21, 121, 43,
		78, 212, 229, 172, 115, 243, 167, 87, 7, 112,
		192, 247, 140, 128, 99, 13, 103, 74, 222, 237,
		49, 197, 254, 24, 227, 165, 153, 119, 38, 184,
		180, 124, 17, 68, 146, 217, 35, 32, 137, 46,
		55, 63, 209, 91, 149, 188, 207, 205, 144, 135,
		151, 178, 220, 252, 190, 97, 242, 86, 211, 171,
		20, 42, 93, 158, 132, 60, 57, 83, 71, 109,
		65, 162, 31, 45, 67, 216, 183, 123, 164, 118,
		196, 23, 73, 236, 127, 12, 111, 246, 108, 161,
		59, 82, 41, 157, 85, 170, 251, 96, 134, 177,
		187, 204, 62, 90, 203, 89, 95, 176, 156, 169,
		160, 81, 11, 245, 22, 235, 122, 117, 44, 215,
		79, 174, 213, 233, 230, 231, 173, 232, 116, 214,
		244, 234, 168, 80, 88, 175,
	}
	// galois指数表
	galoisExpTable = []byte{
		1, 2, 4, 8, 16, 32, 64, 128, 29, 58,
		116, 232, 205, 135, 19, 38, 76, 152, 45, 90,
		180, 117, 234, 201, 143, 3, 6, 12, 24, 48,
		96, 192, 157, 39, 78, 156, 37, 74, 148, 53,
		106, 212, 181, 119, 238, 193, 159, 35, 70, 140,
		5, 10, 20, 40, 80, 160, 93, 186, 105, 210,
		185, 111, 222, 161, 95, 190, 97, 194, 153, 47,
		94, 188, 101, 202, 137, 15, 30, 60, 120, 240,
		253, 231, 211, 187, 107, 214, 177, 127, 254, 225,
		223, 163, 91, 182, 113, 226, 217, 175, 67, 134,
		17, 34, 68, 136, 13, 26, 52, 104, 208, 189,
		103, 206, 129, 31, 62, 124, 248, 237, 199, 147,
		59, 118, 236, 197, 151, 51, 102, 204, 133, 23,
		46, 92, 184, 109, 218, 169, 79, 158, 33, 66,
		132, 21, 42, 84, 168, 77, 154, 41, 82, 164,
		85, 170, 73, 146, 57, 114, 228, 213, 183, 115,
		230, 209, 191, 99, 198, 145, 63, 126, 252, 229,
		215, 179, 123, 246, 241, 255, 227, 219, 171, 75,
		150, 49, 98, 196, 149, 55, 110, 220, 165, 87,
		174, 65, 130, 25, 50, 100, 200, 141, 7, 14,
		28, 56, 112, 224, 221, 167, 83, 166, 81, 162,
		89, 178, 121, 242, 249, 239, 195, 155, 43, 86,
		172, 69, 138, 9, 18, 36, 72, 144, 61, 122,
		244, 245, 247, 243, 251, 235, 203, 139, 11, 22,
		44, 88, 176, 125, 250, 233, 207, 131, 27, 54,
		108, 216, 173, 71, 142, 1,
	}
	// 纠错表
	errorCorrectionTable = [maxVersion][maxLevel]*errorCorrection{
		{
			{TotalBytes: 19, BlockECBytes: 7, Group1Block: 1, Group1BlockBytes: 19, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 16, BlockECBytes: 10, Group1Block: 1, Group1BlockBytes: 16, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 13, BlockECBytes: 13, Group1Block: 1, Group1BlockBytes: 13, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 9, BlockECBytes: 17, Group1Block: 1, Group1BlockBytes: 9, Group2Block: 0, Group2BlockBytes: 0},
		},
		{
			{TotalBytes: 34, BlockECBytes: 10, Group1Block: 1, Group1BlockBytes: 34, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 28, BlockECBytes: 16, Group1Block: 1, Group1BlockBytes: 28, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 22, BlockECBytes: 22, Group1Block: 1, Group1BlockBytes: 22, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 16, BlockECBytes: 28, Group1Block: 1, Group1BlockBytes: 16, Group2Block: 0, Group2BlockBytes: 0},
		},
		{
			{TotalBytes: 55, BlockECBytes: 15, Group1Block: 1, Group1BlockBytes: 55, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 44, BlockECBytes: 26, Group1Block: 1, Group1BlockBytes: 44, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 34, BlockECBytes: 18, Group1Block: 2, Group1BlockBytes: 17, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 26, BlockECBytes: 22, Group1Block: 2, Group1BlockBytes: 13, Group2Block: 0, Group2BlockBytes: 0},
		},
		{
			{TotalBytes: 80, BlockECBytes: 20, Group1Block: 1, Group1BlockBytes: 80, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 64, BlockECBytes: 18, Group1Block: 2, Group1BlockBytes: 32, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 48, BlockECBytes: 26, Group1Block: 2, Group1BlockBytes: 24, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 36, BlockECBytes: 16, Group1Block: 4, Group1BlockBytes: 9, Group2Block: 0, Group2BlockBytes: 0},
		},
		{
			{TotalBytes: 108, BlockECBytes: 26, Group1Block: 1, Group1BlockBytes: 108, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 86, BlockECBytes: 24, Group1Block: 2, Group1BlockBytes: 43, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 62, BlockECBytes: 18, Group1Block: 2, Group1BlockBytes: 15, Group2Block: 2, Group2BlockBytes: 16},
			{TotalBytes: 46, BlockECBytes: 22, Group1Block: 2, Group1BlockBytes: 11, Group2Block: 2, Group2BlockBytes: 12},
		},
		{
			{TotalBytes: 136, BlockECBytes: 18, Group1Block: 2, Group1BlockBytes: 68, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 108, BlockECBytes: 16, Group1Block: 4, Group1BlockBytes: 27, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 76, BlockECBytes: 24, Group1Block: 4, Group1BlockBytes: 19, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 60, BlockECBytes: 28, Group1Block: 4, Group1BlockBytes: 15, Group2Block: 0, Group2BlockBytes: 0},
		},
		{
			{TotalBytes: 156, BlockECBytes: 20, Group1Block: 2, Group1BlockBytes: 78, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 124, BlockECBytes: 18, Group1Block: 4, Group1BlockBytes: 31, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 88, BlockECBytes: 18, Group1Block: 2, Group1BlockBytes: 14, Group2Block: 4, Group2BlockBytes: 15},
			{TotalBytes: 66, BlockECBytes: 26, Group1Block: 4, Group1BlockBytes: 13, Group2Block: 1, Group2BlockBytes: 14},
		},
		{
			{TotalBytes: 194, BlockECBytes: 24, Group1Block: 2, Group1BlockBytes: 97, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 154, BlockECBytes: 22, Group1Block: 2, Group1BlockBytes: 38, Group2Block: 2, Group2BlockBytes: 39},
			{TotalBytes: 110, BlockECBytes: 22, Group1Block: 4, Group1BlockBytes: 18, Group2Block: 2, Group2BlockBytes: 19},
			{TotalBytes: 86, BlockECBytes: 26, Group1Block: 4, Group1BlockBytes: 14, Group2Block: 2, Group2BlockBytes: 15},
		},
		{
			{TotalBytes: 232, BlockECBytes: 30, Group1Block: 2, Group1BlockBytes: 116, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 182, BlockECBytes: 22, Group1Block: 3, Group1BlockBytes: 36, Group2Block: 2, Group2BlockBytes: 37},
			{TotalBytes: 132, BlockECBytes: 20, Group1Block: 4, Group1BlockBytes: 16, Group2Block: 4, Group2BlockBytes: 17},
			{TotalBytes: 100, BlockECBytes: 24, Group1Block: 4, Group1BlockBytes: 12, Group2Block: 4, Group2BlockBytes: 13},
		},
		{
			{TotalBytes: 274, BlockECBytes: 18, Group1Block: 2, Group1BlockBytes: 68, Group2Block: 2, Group2BlockBytes: 69},
			{TotalBytes: 216, BlockECBytes: 26, Group1Block: 4, Group1BlockBytes: 43, Group2Block: 1, Group2BlockBytes: 44},
			{TotalBytes: 154, BlockECBytes: 24, Group1Block: 6, Group1BlockBytes: 19, Group2Block: 2, Group2BlockBytes: 20},
			{TotalBytes: 122, BlockECBytes: 28, Group1Block: 6, Group1BlockBytes: 15, Group2Block: 2, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 324, BlockECBytes: 20, Group1Block: 4, Group1BlockBytes: 81, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 254, BlockECBytes: 30, Group1Block: 1, Group1BlockBytes: 50, Group2Block: 4, Group2BlockBytes: 51},
			{TotalBytes: 180, BlockECBytes: 28, Group1Block: 4, Group1BlockBytes: 22, Group2Block: 4, Group2BlockBytes: 23},
			{TotalBytes: 140, BlockECBytes: 24, Group1Block: 3, Group1BlockBytes: 12, Group2Block: 8, Group2BlockBytes: 13},
		},
		{
			{TotalBytes: 370, BlockECBytes: 24, Group1Block: 2, Group1BlockBytes: 92, Group2Block: 2, Group2BlockBytes: 93},
			{TotalBytes: 290, BlockECBytes: 22, Group1Block: 6, Group1BlockBytes: 36, Group2Block: 2, Group2BlockBytes: 37},
			{TotalBytes: 206, BlockECBytes: 26, Group1Block: 4, Group1BlockBytes: 20, Group2Block: 6, Group2BlockBytes: 21},
			{TotalBytes: 158, BlockECBytes: 28, Group1Block: 7, Group1BlockBytes: 14, Group2Block: 4, Group2BlockBytes: 15},
		},
		{
			{TotalBytes: 428, BlockECBytes: 26, Group1Block: 4, Group1BlockBytes: 107, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 334, BlockECBytes: 22, Group1Block: 8, Group1BlockBytes: 37, Group2Block: 1, Group2BlockBytes: 38},
			{TotalBytes: 244, BlockECBytes: 24, Group1Block: 8, Group1BlockBytes: 20, Group2Block: 4, Group2BlockBytes: 21},
			{TotalBytes: 180, BlockECBytes: 22, Group1Block: 12, Group1BlockBytes: 11, Group2Block: 4, Group2BlockBytes: 12},
		},
		{
			{TotalBytes: 461, BlockECBytes: 30, Group1Block: 3, Group1BlockBytes: 115, Group2Block: 1, Group2BlockBytes: 116},
			{TotalBytes: 365, BlockECBytes: 24, Group1Block: 4, Group1BlockBytes: 40, Group2Block: 5, Group2BlockBytes: 41},
			{TotalBytes: 261, BlockECBytes: 20, Group1Block: 11, Group1BlockBytes: 16, Group2Block: 5, Group2BlockBytes: 17},
			{TotalBytes: 197, BlockECBytes: 24, Group1Block: 11, Group1BlockBytes: 12, Group2Block: 5, Group2BlockBytes: 13},
		},
		{
			{TotalBytes: 523, BlockECBytes: 22, Group1Block: 5, Group1BlockBytes: 87, Group2Block: 1, Group2BlockBytes: 88},
			{TotalBytes: 415, BlockECBytes: 24, Group1Block: 5, Group1BlockBytes: 41, Group2Block: 5, Group2BlockBytes: 42},
			{TotalBytes: 295, BlockECBytes: 30, Group1Block: 5, Group1BlockBytes: 24, Group2Block: 7, Group2BlockBytes: 25},
			{TotalBytes: 223, BlockECBytes: 24, Group1Block: 11, Group1BlockBytes: 12, Group2Block: 7, Group2BlockBytes: 13},
		},
		{
			{TotalBytes: 589, BlockECBytes: 24, Group1Block: 5, Group1BlockBytes: 98, Group2Block: 1, Group2BlockBytes: 99},
			{TotalBytes: 453, BlockECBytes: 28, Group1Block: 7, Group1BlockBytes: 45, Group2Block: 3, Group2BlockBytes: 46},
			{TotalBytes: 325, BlockECBytes: 24, Group1Block: 15, Group1BlockBytes: 19, Group2Block: 2, Group2BlockBytes: 20},
			{TotalBytes: 253, BlockECBytes: 30, Group1Block: 3, Group1BlockBytes: 15, Group2Block: 13, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 647, BlockECBytes: 28, Group1Block: 1, Group1BlockBytes: 107, Group2Block: 5, Group2BlockBytes: 108},
			{TotalBytes: 507, BlockECBytes: 28, Group1Block: 10, Group1BlockBytes: 46, Group2Block: 1, Group2BlockBytes: 47},
			{TotalBytes: 367, BlockECBytes: 28, Group1Block: 1, Group1BlockBytes: 22, Group2Block: 15, Group2BlockBytes: 23},
			{TotalBytes: 283, BlockECBytes: 28, Group1Block: 2, Group1BlockBytes: 14, Group2Block: 17, Group2BlockBytes: 15},
		},
		{
			{TotalBytes: 721, BlockECBytes: 30, Group1Block: 5, Group1BlockBytes: 120, Group2Block: 1, Group2BlockBytes: 121},
			{TotalBytes: 563, BlockECBytes: 26, Group1Block: 9, Group1BlockBytes: 43, Group2Block: 4, Group2BlockBytes: 44},
			{TotalBytes: 397, BlockECBytes: 28, Group1Block: 17, Group1BlockBytes: 22, Group2Block: 1, Group2BlockBytes: 23},
			{TotalBytes: 313, BlockECBytes: 28, Group1Block: 2, Group1BlockBytes: 14, Group2Block: 19, Group2BlockBytes: 15},
		},
		{
			{TotalBytes: 795, BlockECBytes: 28, Group1Block: 3, Group1BlockBytes: 113, Group2Block: 4, Group2BlockBytes: 114},
			{TotalBytes: 627, BlockECBytes: 26, Group1Block: 3, Group1BlockBytes: 44, Group2Block: 11, Group2BlockBytes: 45},
			{TotalBytes: 445, BlockECBytes: 26, Group1Block: 17, Group1BlockBytes: 21, Group2Block: 4, Group2BlockBytes: 22},
			{TotalBytes: 341, BlockECBytes: 26, Group1Block: 9, Group1BlockBytes: 13, Group2Block: 16, Group2BlockBytes: 14},
		},
		{
			{TotalBytes: 861, BlockECBytes: 28, Group1Block: 3, Group1BlockBytes: 107, Group2Block: 5, Group2BlockBytes: 108},
			{TotalBytes: 669, BlockECBytes: 26, Group1Block: 3, Group1BlockBytes: 41, Group2Block: 13, Group2BlockBytes: 42},
			{TotalBytes: 485, BlockECBytes: 30, Group1Block: 15, Group1BlockBytes: 24, Group2Block: 5, Group2BlockBytes: 25},
			{TotalBytes: 385, BlockECBytes: 28, Group1Block: 15, Group1BlockBytes: 15, Group2Block: 10, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 932, BlockECBytes: 28, Group1Block: 4, Group1BlockBytes: 116, Group2Block: 4, Group2BlockBytes: 117},
			{TotalBytes: 714, BlockECBytes: 26, Group1Block: 17, Group1BlockBytes: 42, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 512, BlockECBytes: 28, Group1Block: 17, Group1BlockBytes: 22, Group2Block: 6, Group2BlockBytes: 23},
			{TotalBytes: 406, BlockECBytes: 30, Group1Block: 19, Group1BlockBytes: 16, Group2Block: 6, Group2BlockBytes: 17},
		},
		{
			{TotalBytes: 1006, BlockECBytes: 28, Group1Block: 2, Group1BlockBytes: 111, Group2Block: 7, Group2BlockBytes: 112},
			{TotalBytes: 782, BlockECBytes: 28, Group1Block: 17, Group1BlockBytes: 46, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 568, BlockECBytes: 30, Group1Block: 7, Group1BlockBytes: 24, Group2Block: 16, Group2BlockBytes: 25},
			{TotalBytes: 442, BlockECBytes: 24, Group1Block: 34, Group1BlockBytes: 13, Group2Block: 0, Group2BlockBytes: 0},
		},
		{
			{TotalBytes: 1094, BlockECBytes: 30, Group1Block: 4, Group1BlockBytes: 121, Group2Block: 5, Group2BlockBytes: 122},
			{TotalBytes: 860, BlockECBytes: 28, Group1Block: 4, Group1BlockBytes: 47, Group2Block: 14, Group2BlockBytes: 48},
			{TotalBytes: 614, BlockECBytes: 30, Group1Block: 11, Group1BlockBytes: 24, Group2Block: 14, Group2BlockBytes: 25},
			{TotalBytes: 464, BlockECBytes: 30, Group1Block: 16, Group1BlockBytes: 15, Group2Block: 14, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 1174, BlockECBytes: 30, Group1Block: 6, Group1BlockBytes: 117, Group2Block: 4, Group2BlockBytes: 118},
			{TotalBytes: 914, BlockECBytes: 28, Group1Block: 6, Group1BlockBytes: 45, Group2Block: 14, Group2BlockBytes: 46},
			{TotalBytes: 664, BlockECBytes: 30, Group1Block: 11, Group1BlockBytes: 24, Group2Block: 16, Group2BlockBytes: 25},
			{TotalBytes: 514, BlockECBytes: 30, Group1Block: 30, Group1BlockBytes: 16, Group2Block: 2, Group2BlockBytes: 17},
		},
		{
			{TotalBytes: 1276, BlockECBytes: 26, Group1Block: 8, Group1BlockBytes: 106, Group2Block: 4, Group2BlockBytes: 107},
			{TotalBytes: 1000, BlockECBytes: 28, Group1Block: 8, Group1BlockBytes: 47, Group2Block: 13, Group2BlockBytes: 48},
			{TotalBytes: 718, BlockECBytes: 30, Group1Block: 7, Group1BlockBytes: 24, Group2Block: 22, Group2BlockBytes: 25},
			{TotalBytes: 538, BlockECBytes: 30, Group1Block: 22, Group1BlockBytes: 15, Group2Block: 13, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 1370, BlockECBytes: 28, Group1Block: 10, Group1BlockBytes: 114, Group2Block: 2, Group2BlockBytes: 115},
			{TotalBytes: 1062, BlockECBytes: 28, Group1Block: 19, Group1BlockBytes: 46, Group2Block: 4, Group2BlockBytes: 47},
			{TotalBytes: 754, BlockECBytes: 28, Group1Block: 28, Group1BlockBytes: 22, Group2Block: 6, Group2BlockBytes: 23},
			{TotalBytes: 596, BlockECBytes: 30, Group1Block: 33, Group1BlockBytes: 16, Group2Block: 4, Group2BlockBytes: 17},
		},
		{
			{TotalBytes: 1468, BlockECBytes: 30, Group1Block: 8, Group1BlockBytes: 122, Group2Block: 4, Group2BlockBytes: 123},
			{TotalBytes: 1128, BlockECBytes: 28, Group1Block: 22, Group1BlockBytes: 45, Group2Block: 3, Group2BlockBytes: 46},
			{TotalBytes: 808, BlockECBytes: 30, Group1Block: 8, Group1BlockBytes: 23, Group2Block: 26, Group2BlockBytes: 24},
			{TotalBytes: 628, BlockECBytes: 30, Group1Block: 12, Group1BlockBytes: 15, Group2Block: 28, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 1531, BlockECBytes: 30, Group1Block: 3, Group1BlockBytes: 117, Group2Block: 10, Group2BlockBytes: 118},
			{TotalBytes: 1193, BlockECBytes: 28, Group1Block: 3, Group1BlockBytes: 45, Group2Block: 23, Group2BlockBytes: 46},
			{TotalBytes: 871, BlockECBytes: 30, Group1Block: 4, Group1BlockBytes: 24, Group2Block: 31, Group2BlockBytes: 25},
			{TotalBytes: 661, BlockECBytes: 30, Group1Block: 11, Group1BlockBytes: 15, Group2Block: 31, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 1631, BlockECBytes: 30, Group1Block: 7, Group1BlockBytes: 116, Group2Block: 7, Group2BlockBytes: 117},
			{TotalBytes: 1267, BlockECBytes: 28, Group1Block: 21, Group1BlockBytes: 45, Group2Block: 7, Group2BlockBytes: 46},
			{TotalBytes: 911, BlockECBytes: 30, Group1Block: 1, Group1BlockBytes: 23, Group2Block: 37, Group2BlockBytes: 24},
			{TotalBytes: 701, BlockECBytes: 30, Group1Block: 19, Group1BlockBytes: 15, Group2Block: 26, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 1735, BlockECBytes: 30, Group1Block: 5, Group1BlockBytes: 115, Group2Block: 10, Group2BlockBytes: 116},
			{TotalBytes: 1373, BlockECBytes: 28, Group1Block: 19, Group1BlockBytes: 47, Group2Block: 10, Group2BlockBytes: 48},
			{TotalBytes: 985, BlockECBytes: 30, Group1Block: 15, Group1BlockBytes: 24, Group2Block: 25, Group2BlockBytes: 25},
			{TotalBytes: 745, BlockECBytes: 30, Group1Block: 23, Group1BlockBytes: 15, Group2Block: 25, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 1843, BlockECBytes: 30, Group1Block: 13, Group1BlockBytes: 115, Group2Block: 3, Group2BlockBytes: 116},
			{TotalBytes: 1455, BlockECBytes: 28, Group1Block: 2, Group1BlockBytes: 46, Group2Block: 29, Group2BlockBytes: 47},
			{TotalBytes: 1033, BlockECBytes: 30, Group1Block: 42, Group1BlockBytes: 24, Group2Block: 1, Group2BlockBytes: 25},
			{TotalBytes: 793, BlockECBytes: 30, Group1Block: 23, Group1BlockBytes: 15, Group2Block: 28, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 1955, BlockECBytes: 30, Group1Block: 17, Group1BlockBytes: 115, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 1541, BlockECBytes: 28, Group1Block: 10, Group1BlockBytes: 46, Group2Block: 23, Group2BlockBytes: 47},
			{TotalBytes: 1115, BlockECBytes: 30, Group1Block: 10, Group1BlockBytes: 24, Group2Block: 35, Group2BlockBytes: 25},
			{TotalBytes: 845, BlockECBytes: 30, Group1Block: 19, Group1BlockBytes: 15, Group2Block: 35, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 2071, BlockECBytes: 30, Group1Block: 17, Group1BlockBytes: 115, Group2Block: 1, Group2BlockBytes: 116},
			{TotalBytes: 1631, BlockECBytes: 28, Group1Block: 14, Group1BlockBytes: 46, Group2Block: 21, Group2BlockBytes: 47},
			{TotalBytes: 1171, BlockECBytes: 30, Group1Block: 29, Group1BlockBytes: 24, Group2Block: 19, Group2BlockBytes: 25},
			{TotalBytes: 901, BlockECBytes: 30, Group1Block: 11, Group1BlockBytes: 15, Group2Block: 46, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 2191, BlockECBytes: 30, Group1Block: 13, Group1BlockBytes: 115, Group2Block: 6, Group2BlockBytes: 116},
			{TotalBytes: 1725, BlockECBytes: 28, Group1Block: 14, Group1BlockBytes: 46, Group2Block: 23, Group2BlockBytes: 47},
			{TotalBytes: 1231, BlockECBytes: 30, Group1Block: 44, Group1BlockBytes: 24, Group2Block: 7, Group2BlockBytes: 25},
			{TotalBytes: 961, BlockECBytes: 30, Group1Block: 59, Group1BlockBytes: 16, Group2Block: 1, Group2BlockBytes: 17},
		},
		{
			{TotalBytes: 2306, BlockECBytes: 30, Group1Block: 12, Group1BlockBytes: 121, Group2Block: 7, Group2BlockBytes: 122},
			{TotalBytes: 1812, BlockECBytes: 28, Group1Block: 12, Group1BlockBytes: 47, Group2Block: 26, Group2BlockBytes: 48},
			{TotalBytes: 1286, BlockECBytes: 30, Group1Block: 39, Group1BlockBytes: 24, Group2Block: 14, Group2BlockBytes: 25},
			{TotalBytes: 986, BlockECBytes: 30, Group1Block: 22, Group1BlockBytes: 15, Group2Block: 41, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 2434, BlockECBytes: 30, Group1Block: 6, Group1BlockBytes: 121, Group2Block: 14, Group2BlockBytes: 122},
			{TotalBytes: 1914, BlockECBytes: 28, Group1Block: 6, Group1BlockBytes: 47, Group2Block: 34, Group2BlockBytes: 48},
			{TotalBytes: 1354, BlockECBytes: 30, Group1Block: 46, Group1BlockBytes: 24, Group2Block: 10, Group2BlockBytes: 25},
			{TotalBytes: 1054, BlockECBytes: 30, Group1Block: 2, Group1BlockBytes: 15, Group2Block: 64, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 2566, BlockECBytes: 30, Group1Block: 17, Group1BlockBytes: 122, Group2Block: 4, Group2BlockBytes: 123},
			{TotalBytes: 1992, BlockECBytes: 28, Group1Block: 29, Group1BlockBytes: 46, Group2Block: 14, Group2BlockBytes: 47},
			{TotalBytes: 1426, BlockECBytes: 30, Group1Block: 49, Group1BlockBytes: 24, Group2Block: 10, Group2BlockBytes: 25},
			{TotalBytes: 1096, BlockECBytes: 30, Group1Block: 24, Group1BlockBytes: 15, Group2Block: 46, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 2702, BlockECBytes: 30, Group1Block: 4, Group1BlockBytes: 122, Group2Block: 18, Group2BlockBytes: 123},
			{TotalBytes: 2102, BlockECBytes: 28, Group1Block: 13, Group1BlockBytes: 46, Group2Block: 32, Group2BlockBytes: 47},
			{TotalBytes: 1502, BlockECBytes: 30, Group1Block: 48, Group1BlockBytes: 24, Group2Block: 14, Group2BlockBytes: 25},
			{TotalBytes: 1142, BlockECBytes: 30, Group1Block: 42, Group1BlockBytes: 15, Group2Block: 32, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 2812, BlockECBytes: 30, Group1Block: 20, Group1BlockBytes: 117, Group2Block: 4, Group2BlockBytes: 118},
			{TotalBytes: 2216, BlockECBytes: 28, Group1Block: 40, Group1BlockBytes: 47, Group2Block: 7, Group2BlockBytes: 48},
			{TotalBytes: 1582, BlockECBytes: 30, Group1Block: 43, Group1BlockBytes: 24, Group2Block: 22, Group2BlockBytes: 25},
			{TotalBytes: 1222, BlockECBytes: 30, Group1Block: 10, Group1BlockBytes: 15, Group2Block: 67, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 2956, BlockECBytes: 30, Group1Block: 19, Group1BlockBytes: 118, Group2Block: 6, Group2BlockBytes: 119},
			{TotalBytes: 2334, BlockECBytes: 28, Group1Block: 18, Group1BlockBytes: 47, Group2Block: 31, Group2BlockBytes: 48},
			{TotalBytes: 1666, BlockECBytes: 30, Group1Block: 34, Group1BlockBytes: 24, Group2Block: 34, Group2BlockBytes: 25},
			{TotalBytes: 1276, BlockECBytes: 30, Group1Block: 20, Group1BlockBytes: 15, Group2Block: 61, Group2BlockBytes: 16},
		},
	}
)

// 纠错表
type errorCorrection struct {
	TotalBytes       int // Total Number of Data Codewords for this Version and EC Level
	BlockECBytes     int // EC Codewords Per Block
	Group1Block      int // Number of Blocks in Group 1
	Group1BlockBytes int // Number of Data Codewords in Each of Group 1's Blocks
	Group2Block      int // Number of Blocks in Group 2
	Group2BlockBytes int // Number of Data Codewords in Each of Group 2's Blocks
}

type eccEncoder struct {
	buff *buffer  // 共享的缓存
	poly []byte   // 生成的多项式
	data []byte   // 编码后的数据
	xy   [][]byte // 二维表，交错使用
}

// 对b进行编码，并返回编码后的数据
func (e *eccEncoder) Encode(data []byte, version version, level Level) {
	// 纠错表
	ec := errorCorrectionTable[version][level]
	// 生成多项式
	e.genPoly(ec.BlockECBytes)
	// 编码
	e.data = e.data[:0]
	e.data = append(e.data, data...)
	p := data
	for i := 0; i < ec.Group1Block; i++ {
		e.data = append(e.data, e.encode(p[:ec.Group1BlockBytes])...)
		p = p[ec.Group1BlockBytes:]
	}
	for i := 0; i < ec.Group2Block; i++ {
		e.data = append(e.data, e.encode(p[:ec.Group2BlockBytes])...)
		p = p[ec.Group2BlockBytes:]
	}
	// 交错
	if ec.Group2Block > 0 || ec.Group1Block > 1 {
		e.buff.Resize(len(e.data), -1)
		// 二维表
		e.xy = e.xy[:0]
		p = data
		for i := 0; i < ec.Group1Block; i++ {
			e.xy = append(e.xy, p[:ec.Group1BlockBytes])
			p = p[ec.Group1BlockBytes:]
		}
		for i := 0; i < ec.Group2Block; i++ {
			e.xy = append(e.xy, p[:ec.Group2BlockBytes])
			p = p[ec.Group2BlockBytes:]
		}
		idx := 0
		// 数据交错
		for x := 0; x < ec.Group1BlockBytes; x++ {
			for y := 0; y < len(e.xy); y++ {
				e.buff.data[idx] = e.xy[y][x]
				idx++
			}
		}
		if ec.Group2BlockBytes > ec.Group1BlockBytes {
			for y := 0; y < ec.Group2Block; y++ {
				e.buff.data[idx] = e.xy[y+ec.Group1Block][ec.Group1BlockBytes]
				idx++
			}
		}
		// 纠错码交错
		e.xy = e.xy[:0]
		p = e.data[len(data):]
		for i := 0; i < (ec.Group1Block + ec.Group2Block); i++ {
			e.xy = append(e.xy, p[:ec.BlockECBytes])
			p = data[ec.BlockECBytes:]
		}
		for i := 0; i < ec.Group2Block; i++ {
			e.xy = append(e.xy, p[:ec.BlockECBytes])
			p = p[ec.BlockECBytes:]
		}
		for x := 0; x < ec.BlockECBytes; x++ {
			for y := 0; y < len(e.xy); y++ {
				e.buff.data[idx] = e.xy[y][x]
				idx++
			}
		}
		// 添加余数
		if interleaveRemainder[level] > 0 {
			e.buff.data = append(e.buff.data, 0)
		}
		// 交换缓存
		t := e.buff.data
		e.buff.data = e.data
		e.data = t
	}
}

// 编码
func (e *eccEncoder) encode(data []byte) []byte {
	e.buff.Resize(len(data)+len(e.poly)-1, len(data))
	copy(e.buff.data, data)
	for i := 0; i < len(data); i++ {
		if e.buff.data[i] != 0 {
			for j := 1; j < len(e.poly); j++ {
				e.buff.data[i+j] ^= e.galoisMul(e.poly[j], e.buff.data[i])
			}
		}
	}
	return e.buff.data[len(data):]
}

// 多项式乘法，p.poly * poly
func (e *eccEncoder) mulPoly(poly []byte) {
	e.buff.Resize(len(e.poly)+len(poly)-1, 0)
	for i := 0; i < len(poly); i++ {
		for j := 0; j < len(e.poly); j++ {
			e.buff.data[i+j] ^= e.galoisMul(e.poly[j], poly[i])
		}
	}
	// 交换缓存
	t := e.buff.data
	e.buff.data = e.poly
	e.poly = t
}

// 根据n，生成多项式
func (e *eccEncoder) genPoly(n int) {
	e.poly = e.poly[:1]
	e.poly[0] = 1
	a := []byte{1, 0}
	for i := 0; i < n; i++ {
		a[1] = galoisExpTable[i]
		e.mulPoly(a[:])
	}
}

// galois两个数乘法
func (e *eccEncoder) galoisMul(n1, n2 byte) byte {
	if n1 == 0 || n2 == 0 {
		return 0
	}
	return galoisExpTable[(int(galoisLogTable[n1])+int(galoisLogTable[n2]))%255]
}
