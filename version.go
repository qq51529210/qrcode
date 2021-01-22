package qrcode

import "fmt"

const (
	maxVersion = 40
)

var (
	versionCapacity = [maxLevel][maxMode][maxVersion]int{
		{
			{
				41, 77, 127, 187, 255, 322, 370, 461, 552, 652, 772,
				883, 1022, 1101, 1250, 1408, 1548, 1725, 1903, 2061, 2232,
				2409, 2620, 2812, 3057, 3283, 3517, 3669, 3909, 4158, 4417,
				4686, 4965, 5253, 5529, 5836, 6153, 6479, 6743, 7089,
			},
			{
				25, 47, 77, 114, 154, 195, 224, 279, 335, 395, 468,
				535, 619, 667, 758, 854, 938, 1046, 1153, 1249, 1352,
				1460, 1588, 1704, 1853, 1990, 2132, 2223, 2369, 2520, 2677,
				2840, 3009, 3183, 3351, 3537, 3729, 3927, 4087, 4296,
			},
			{
				17, 32, 53, 78, 106, 134, 154, 192, 230, 271, 321,
				367, 425, 458, 520, 586, 644, 718, 792, 858, 929,
				1003, 1091, 1171, 1273, 1367, 1465, 1528, 1628, 1732, 1840,
				1952, 2068, 2188, 2303, 2431, 2563, 2699, 2809, 2953,
			},
			{
				10, 20, 32, 48, 65, 82, 95, 118, 141, 167, 198,
				226, 262, 282, 320, 361, 397, 442, 488, 528, 572,
				618, 672, 721, 784, 842, 902, 940, 1002, 1066, 1132,
				1201, 1273, 1347, 1417, 1496, 1577, 1661, 1729, 1817,
			},
		},
		{
			{
				34, 63, 101, 149, 202, 255, 293, 365, 432, 513, 604,
				691, 796, 871, 991, 1082, 1212, 1346, 1500, 1600, 1708,
				1872, 2059, 2188, 2395, 2544, 2701, 2857, 3035, 3289, 3486,
				3693, 3909, 4134, 4343, 4588, 4775, 5039, 5313, 5596,
			},
			{
				20, 38, 61, 90, 122, 154, 178, 221, 262, 311, 366,
				419, 483, 528, 600, 656, 734, 816, 909, 970, 1035,
				1134, 1248, 1326, 1451, 1542, 1637, 1732, 1839, 1994, 2113,
				2238, 2369, 2506, 2632, 2780, 2894, 3054, 3220, 3391,
			},
			{
				20, 38, 61, 90, 122, 154, 178, 221, 262, 311, 366,
				419, 483, 528, 600, 656, 734, 816, 909, 970, 1035,
				1134, 1248, 1326, 1451, 1542, 1637, 1732, 1839, 1994, 2113,
				2238, 2369, 2506, 2632, 2780, 2894, 3054, 3220, 3391,
			},
			{
				8, 16, 26, 38, 52, 65, 75, 93, 111, 131, 155,
				177, 204, 223, 254, 277, 310, 345, 384, 410, 438,
				480, 528, 561, 614, 652, 692, 732, 778, 843, 894,
				947, 1002, 1060, 1113, 1176, 1224, 1292, 1362, 1435,
			},
		},
		{
			{
				27, 48, 77, 111, 144, 178, 207, 259, 312, 364, 427,
				489, 580, 621, 703, 775, 876, 948, 1063, 1159, 1224,
				1358, 1468, 1588, 1718, 1804, 1933, 2085, 2181, 2358, 2473,
				2670, 2805, 2949, 3081, 3244, 3417, 3599, 3791, 3993,
			},
			{
				16, 29, 47, 67, 87, 108, 125, 157, 189, 221, 259,
				296, 352, 376, 426, 470, 531, 574, 644, 702, 742,
				823, 890, 963, 1041, 1094, 1172, 1263, 1322, 1429, 1499,
				1618, 1700, 1787, 1867, 1966, 2071, 2181, 2298, 2420,
			},
			{
				11, 20, 32, 46, 60, 74, 86, 108, 130, 151, 177,
				203, 241, 258, 292, 322, 364, 394, 442, 482, 509,
				565, 611, 661, 715, 751, 805, 868, 908, 982, 1030,
				1112, 1168, 1228, 1283, 1351, 1423, 1499, 1579, 1663,
			},
			{
				7, 12, 20, 28, 37, 45, 53, 66, 80, 93, 109,
				125, 149, 159, 180, 198, 224, 243, 272, 297, 314,
				348, 376, 407, 440, 462, 496, 534, 559, 604, 634,
				684, 719, 756, 790, 832, 876, 923, 972, 1024,
			},
		},
		{
			{
				17, 34, 58, 82, 106, 139, 154, 202, 235, 288, 331,
				374, 427, 468, 530, 602, 674, 746, 813, 919, 969,
				1056, 1108, 1228, 1286, 1425, 1501, 1581, 1677, 1782, 1897,
				2022, 2157, 2301, 2361, 2524, 2625, 2735, 2927, 3057,
			},
			{
				10, 20, 35, 50, 64, 84, 93, 122, 143, 174, 200,
				227, 259, 283, 321, 365, 408, 452, 493, 557, 587,
				640, 672, 744, 779, 864, 910, 958, 1016, 1080, 1150,
				1226, 1307, 1394, 1431, 1530, 1591, 1658, 1774, 1852,
			},
			{
				7, 14, 24, 34, 44, 58, 64, 84, 98, 119, 137,
				155, 177, 194, 220, 250, 280, 310, 338, 382, 403,
				439, 461, 511, 535, 593, 625, 658, 698, 742, 790,
				842, 898, 958, 983, 1051, 1093, 1139, 1219, 1273,
			},
			{
				4, 8, 15, 21, 27, 36, 39, 52, 60, 74, 85,
				96, 109, 120, 136, 154, 173, 191, 208, 235, 248,
				270, 284, 315, 330, 365, 385, 405, 430, 457, 486,
				518, 553, 590, 605, 647, 673, 701, 750, 784,
			},
		},
	} // 用于快速判断版本
	versionECTable = [maxVersion][maxLevel]*ecTable{
		{
			{TotalBytes: 19, BlockBytes: 7, Group1Block: 1, Group1BlockBytes: 19, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 16, BlockBytes: 10, Group1Block: 1, Group1BlockBytes: 16, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 13, BlockBytes: 13, Group1Block: 1, Group1BlockBytes: 13, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 9, BlockBytes: 17, Group1Block: 1, Group1BlockBytes: 9, Group2Block: 0, Group2BlockBytes: 0},
		},
		{
			{TotalBytes: 34, BlockBytes: 10, Group1Block: 1, Group1BlockBytes: 34, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 28, BlockBytes: 16, Group1Block: 1, Group1BlockBytes: 28, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 22, BlockBytes: 22, Group1Block: 1, Group1BlockBytes: 22, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 16, BlockBytes: 28, Group1Block: 1, Group1BlockBytes: 16, Group2Block: 0, Group2BlockBytes: 0},
		},
		{
			{TotalBytes: 55, BlockBytes: 15, Group1Block: 1, Group1BlockBytes: 55, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 44, BlockBytes: 26, Group1Block: 1, Group1BlockBytes: 44, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 34, BlockBytes: 18, Group1Block: 2, Group1BlockBytes: 17, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 26, BlockBytes: 22, Group1Block: 2, Group1BlockBytes: 13, Group2Block: 0, Group2BlockBytes: 0},
		},
		{
			{TotalBytes: 80, BlockBytes: 20, Group1Block: 1, Group1BlockBytes: 80, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 64, BlockBytes: 18, Group1Block: 2, Group1BlockBytes: 32, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 48, BlockBytes: 26, Group1Block: 2, Group1BlockBytes: 24, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 36, BlockBytes: 16, Group1Block: 4, Group1BlockBytes: 9, Group2Block: 0, Group2BlockBytes: 0},
		},
		{
			{TotalBytes: 108, BlockBytes: 26, Group1Block: 1, Group1BlockBytes: 108, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 86, BlockBytes: 24, Group1Block: 2, Group1BlockBytes: 43, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 62, BlockBytes: 18, Group1Block: 2, Group1BlockBytes: 15, Group2Block: 2, Group2BlockBytes: 16},
			{TotalBytes: 46, BlockBytes: 22, Group1Block: 2, Group1BlockBytes: 11, Group2Block: 2, Group2BlockBytes: 12},
		},
		{
			{TotalBytes: 136, BlockBytes: 18, Group1Block: 2, Group1BlockBytes: 68, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 108, BlockBytes: 16, Group1Block: 4, Group1BlockBytes: 27, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 76, BlockBytes: 24, Group1Block: 4, Group1BlockBytes: 19, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 60, BlockBytes: 28, Group1Block: 4, Group1BlockBytes: 15, Group2Block: 0, Group2BlockBytes: 0},
		},
		{
			{TotalBytes: 156, BlockBytes: 20, Group1Block: 2, Group1BlockBytes: 78, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 124, BlockBytes: 18, Group1Block: 4, Group1BlockBytes: 31, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 88, BlockBytes: 18, Group1Block: 2, Group1BlockBytes: 14, Group2Block: 4, Group2BlockBytes: 15},
			{TotalBytes: 66, BlockBytes: 26, Group1Block: 4, Group1BlockBytes: 13, Group2Block: 1, Group2BlockBytes: 14},
		},
		{
			{TotalBytes: 194, BlockBytes: 24, Group1Block: 2, Group1BlockBytes: 97, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 154, BlockBytes: 22, Group1Block: 2, Group1BlockBytes: 38, Group2Block: 2, Group2BlockBytes: 39},
			{TotalBytes: 110, BlockBytes: 22, Group1Block: 4, Group1BlockBytes: 18, Group2Block: 2, Group2BlockBytes: 19},
			{TotalBytes: 86, BlockBytes: 26, Group1Block: 4, Group1BlockBytes: 14, Group2Block: 2, Group2BlockBytes: 15},
		},
		{
			{TotalBytes: 232, BlockBytes: 30, Group1Block: 2, Group1BlockBytes: 116, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 182, BlockBytes: 22, Group1Block: 3, Group1BlockBytes: 36, Group2Block: 2, Group2BlockBytes: 37},
			{TotalBytes: 132, BlockBytes: 20, Group1Block: 4, Group1BlockBytes: 16, Group2Block: 4, Group2BlockBytes: 17},
			{TotalBytes: 100, BlockBytes: 24, Group1Block: 4, Group1BlockBytes: 12, Group2Block: 4, Group2BlockBytes: 13},
		},
		{
			{TotalBytes: 274, BlockBytes: 18, Group1Block: 2, Group1BlockBytes: 68, Group2Block: 2, Group2BlockBytes: 69},
			{TotalBytes: 216, BlockBytes: 26, Group1Block: 4, Group1BlockBytes: 43, Group2Block: 1, Group2BlockBytes: 44},
			{TotalBytes: 154, BlockBytes: 24, Group1Block: 6, Group1BlockBytes: 19, Group2Block: 2, Group2BlockBytes: 20},
			{TotalBytes: 122, BlockBytes: 28, Group1Block: 6, Group1BlockBytes: 15, Group2Block: 2, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 324, BlockBytes: 20, Group1Block: 4, Group1BlockBytes: 81, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 254, BlockBytes: 30, Group1Block: 1, Group1BlockBytes: 50, Group2Block: 4, Group2BlockBytes: 51},
			{TotalBytes: 180, BlockBytes: 28, Group1Block: 4, Group1BlockBytes: 22, Group2Block: 4, Group2BlockBytes: 23},
			{TotalBytes: 140, BlockBytes: 24, Group1Block: 3, Group1BlockBytes: 12, Group2Block: 8, Group2BlockBytes: 13},
		},
		{
			{TotalBytes: 370, BlockBytes: 24, Group1Block: 2, Group1BlockBytes: 92, Group2Block: 2, Group2BlockBytes: 93},
			{TotalBytes: 290, BlockBytes: 22, Group1Block: 6, Group1BlockBytes: 36, Group2Block: 2, Group2BlockBytes: 37},
			{TotalBytes: 206, BlockBytes: 26, Group1Block: 4, Group1BlockBytes: 20, Group2Block: 6, Group2BlockBytes: 21},
			{TotalBytes: 158, BlockBytes: 28, Group1Block: 7, Group1BlockBytes: 14, Group2Block: 4, Group2BlockBytes: 15},
		},
		{
			{TotalBytes: 428, BlockBytes: 26, Group1Block: 4, Group1BlockBytes: 107, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 334, BlockBytes: 22, Group1Block: 8, Group1BlockBytes: 37, Group2Block: 1, Group2BlockBytes: 38},
			{TotalBytes: 244, BlockBytes: 24, Group1Block: 8, Group1BlockBytes: 20, Group2Block: 4, Group2BlockBytes: 21},
			{TotalBytes: 180, BlockBytes: 22, Group1Block: 12, Group1BlockBytes: 11, Group2Block: 4, Group2BlockBytes: 12},
		},
		{
			{TotalBytes: 461, BlockBytes: 30, Group1Block: 3, Group1BlockBytes: 115, Group2Block: 1, Group2BlockBytes: 116},
			{TotalBytes: 365, BlockBytes: 24, Group1Block: 4, Group1BlockBytes: 40, Group2Block: 5, Group2BlockBytes: 41},
			{TotalBytes: 261, BlockBytes: 20, Group1Block: 11, Group1BlockBytes: 16, Group2Block: 5, Group2BlockBytes: 17},
			{TotalBytes: 197, BlockBytes: 24, Group1Block: 11, Group1BlockBytes: 12, Group2Block: 5, Group2BlockBytes: 13},
		},
		{
			{TotalBytes: 523, BlockBytes: 22, Group1Block: 5, Group1BlockBytes: 87, Group2Block: 1, Group2BlockBytes: 88},
			{TotalBytes: 415, BlockBytes: 24, Group1Block: 5, Group1BlockBytes: 41, Group2Block: 5, Group2BlockBytes: 42},
			{TotalBytes: 295, BlockBytes: 30, Group1Block: 5, Group1BlockBytes: 24, Group2Block: 7, Group2BlockBytes: 25},
			{TotalBytes: 223, BlockBytes: 24, Group1Block: 11, Group1BlockBytes: 12, Group2Block: 7, Group2BlockBytes: 13},
		},
		{
			{TotalBytes: 589, BlockBytes: 24, Group1Block: 5, Group1BlockBytes: 98, Group2Block: 1, Group2BlockBytes: 99},
			{TotalBytes: 453, BlockBytes: 28, Group1Block: 7, Group1BlockBytes: 45, Group2Block: 3, Group2BlockBytes: 46},
			{TotalBytes: 325, BlockBytes: 24, Group1Block: 15, Group1BlockBytes: 19, Group2Block: 2, Group2BlockBytes: 20},
			{TotalBytes: 253, BlockBytes: 30, Group1Block: 3, Group1BlockBytes: 15, Group2Block: 13, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 647, BlockBytes: 28, Group1Block: 1, Group1BlockBytes: 107, Group2Block: 5, Group2BlockBytes: 108},
			{TotalBytes: 507, BlockBytes: 28, Group1Block: 10, Group1BlockBytes: 46, Group2Block: 1, Group2BlockBytes: 47},
			{TotalBytes: 367, BlockBytes: 28, Group1Block: 1, Group1BlockBytes: 22, Group2Block: 15, Group2BlockBytes: 23},
			{TotalBytes: 283, BlockBytes: 28, Group1Block: 2, Group1BlockBytes: 14, Group2Block: 17, Group2BlockBytes: 15},
		},
		{
			{TotalBytes: 721, BlockBytes: 30, Group1Block: 5, Group1BlockBytes: 120, Group2Block: 1, Group2BlockBytes: 121},
			{TotalBytes: 563, BlockBytes: 26, Group1Block: 9, Group1BlockBytes: 43, Group2Block: 4, Group2BlockBytes: 44},
			{TotalBytes: 397, BlockBytes: 28, Group1Block: 17, Group1BlockBytes: 22, Group2Block: 1, Group2BlockBytes: 23},
			{TotalBytes: 313, BlockBytes: 28, Group1Block: 2, Group1BlockBytes: 14, Group2Block: 19, Group2BlockBytes: 15},
		},
		{
			{TotalBytes: 795, BlockBytes: 28, Group1Block: 3, Group1BlockBytes: 113, Group2Block: 4, Group2BlockBytes: 114},
			{TotalBytes: 627, BlockBytes: 26, Group1Block: 3, Group1BlockBytes: 44, Group2Block: 11, Group2BlockBytes: 45},
			{TotalBytes: 445, BlockBytes: 26, Group1Block: 17, Group1BlockBytes: 21, Group2Block: 4, Group2BlockBytes: 22},
			{TotalBytes: 341, BlockBytes: 26, Group1Block: 9, Group1BlockBytes: 13, Group2Block: 16, Group2BlockBytes: 14},
		},
		{
			{TotalBytes: 861, BlockBytes: 28, Group1Block: 3, Group1BlockBytes: 107, Group2Block: 5, Group2BlockBytes: 108},
			{TotalBytes: 669, BlockBytes: 26, Group1Block: 3, Group1BlockBytes: 41, Group2Block: 13, Group2BlockBytes: 42},
			{TotalBytes: 485, BlockBytes: 30, Group1Block: 15, Group1BlockBytes: 24, Group2Block: 5, Group2BlockBytes: 25},
			{TotalBytes: 385, BlockBytes: 28, Group1Block: 15, Group1BlockBytes: 15, Group2Block: 10, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 932, BlockBytes: 28, Group1Block: 4, Group1BlockBytes: 116, Group2Block: 4, Group2BlockBytes: 117},
			{TotalBytes: 714, BlockBytes: 26, Group1Block: 17, Group1BlockBytes: 42, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 512, BlockBytes: 28, Group1Block: 17, Group1BlockBytes: 22, Group2Block: 6, Group2BlockBytes: 23},
			{TotalBytes: 406, BlockBytes: 30, Group1Block: 19, Group1BlockBytes: 16, Group2Block: 6, Group2BlockBytes: 17},
		},
		{
			{TotalBytes: 1006, BlockBytes: 28, Group1Block: 2, Group1BlockBytes: 111, Group2Block: 7, Group2BlockBytes: 112},
			{TotalBytes: 782, BlockBytes: 28, Group1Block: 17, Group1BlockBytes: 46, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 568, BlockBytes: 30, Group1Block: 7, Group1BlockBytes: 24, Group2Block: 16, Group2BlockBytes: 25},
			{TotalBytes: 442, BlockBytes: 24, Group1Block: 34, Group1BlockBytes: 13, Group2Block: 0, Group2BlockBytes: 0},
		},
		{
			{TotalBytes: 1094, BlockBytes: 30, Group1Block: 4, Group1BlockBytes: 121, Group2Block: 5, Group2BlockBytes: 122},
			{TotalBytes: 860, BlockBytes: 28, Group1Block: 4, Group1BlockBytes: 47, Group2Block: 14, Group2BlockBytes: 48},
			{TotalBytes: 614, BlockBytes: 30, Group1Block: 11, Group1BlockBytes: 24, Group2Block: 14, Group2BlockBytes: 25},
			{TotalBytes: 464, BlockBytes: 30, Group1Block: 16, Group1BlockBytes: 15, Group2Block: 14, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 1174, BlockBytes: 30, Group1Block: 6, Group1BlockBytes: 117, Group2Block: 4, Group2BlockBytes: 118},
			{TotalBytes: 914, BlockBytes: 28, Group1Block: 6, Group1BlockBytes: 45, Group2Block: 14, Group2BlockBytes: 46},
			{TotalBytes: 664, BlockBytes: 30, Group1Block: 11, Group1BlockBytes: 24, Group2Block: 16, Group2BlockBytes: 25},
			{TotalBytes: 514, BlockBytes: 30, Group1Block: 30, Group1BlockBytes: 16, Group2Block: 2, Group2BlockBytes: 17},
		},
		{
			{TotalBytes: 1276, BlockBytes: 26, Group1Block: 8, Group1BlockBytes: 106, Group2Block: 4, Group2BlockBytes: 107},
			{TotalBytes: 1000, BlockBytes: 28, Group1Block: 8, Group1BlockBytes: 47, Group2Block: 13, Group2BlockBytes: 48},
			{TotalBytes: 718, BlockBytes: 30, Group1Block: 7, Group1BlockBytes: 24, Group2Block: 22, Group2BlockBytes: 25},
			{TotalBytes: 538, BlockBytes: 30, Group1Block: 22, Group1BlockBytes: 15, Group2Block: 13, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 1370, BlockBytes: 28, Group1Block: 10, Group1BlockBytes: 114, Group2Block: 2, Group2BlockBytes: 115},
			{TotalBytes: 1062, BlockBytes: 28, Group1Block: 19, Group1BlockBytes: 46, Group2Block: 4, Group2BlockBytes: 47},
			{TotalBytes: 754, BlockBytes: 28, Group1Block: 28, Group1BlockBytes: 22, Group2Block: 6, Group2BlockBytes: 23},
			{TotalBytes: 596, BlockBytes: 30, Group1Block: 33, Group1BlockBytes: 16, Group2Block: 4, Group2BlockBytes: 17},
		},
		{
			{TotalBytes: 1468, BlockBytes: 30, Group1Block: 8, Group1BlockBytes: 122, Group2Block: 4, Group2BlockBytes: 123},
			{TotalBytes: 1128, BlockBytes: 28, Group1Block: 22, Group1BlockBytes: 45, Group2Block: 3, Group2BlockBytes: 46},
			{TotalBytes: 808, BlockBytes: 30, Group1Block: 8, Group1BlockBytes: 23, Group2Block: 26, Group2BlockBytes: 24},
			{TotalBytes: 628, BlockBytes: 30, Group1Block: 12, Group1BlockBytes: 15, Group2Block: 28, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 1531, BlockBytes: 30, Group1Block: 3, Group1BlockBytes: 117, Group2Block: 10, Group2BlockBytes: 118},
			{TotalBytes: 1193, BlockBytes: 28, Group1Block: 3, Group1BlockBytes: 45, Group2Block: 23, Group2BlockBytes: 46},
			{TotalBytes: 871, BlockBytes: 30, Group1Block: 4, Group1BlockBytes: 24, Group2Block: 31, Group2BlockBytes: 25},
			{TotalBytes: 661, BlockBytes: 30, Group1Block: 11, Group1BlockBytes: 15, Group2Block: 31, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 1631, BlockBytes: 30, Group1Block: 7, Group1BlockBytes: 116, Group2Block: 7, Group2BlockBytes: 117},
			{TotalBytes: 1267, BlockBytes: 28, Group1Block: 21, Group1BlockBytes: 45, Group2Block: 7, Group2BlockBytes: 46},
			{TotalBytes: 911, BlockBytes: 30, Group1Block: 1, Group1BlockBytes: 23, Group2Block: 37, Group2BlockBytes: 24},
			{TotalBytes: 701, BlockBytes: 30, Group1Block: 19, Group1BlockBytes: 15, Group2Block: 26, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 1735, BlockBytes: 30, Group1Block: 5, Group1BlockBytes: 115, Group2Block: 10, Group2BlockBytes: 116},
			{TotalBytes: 1373, BlockBytes: 28, Group1Block: 19, Group1BlockBytes: 47, Group2Block: 10, Group2BlockBytes: 48},
			{TotalBytes: 985, BlockBytes: 30, Group1Block: 15, Group1BlockBytes: 24, Group2Block: 25, Group2BlockBytes: 25},
			{TotalBytes: 745, BlockBytes: 30, Group1Block: 23, Group1BlockBytes: 15, Group2Block: 25, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 1843, BlockBytes: 30, Group1Block: 13, Group1BlockBytes: 115, Group2Block: 3, Group2BlockBytes: 116},
			{TotalBytes: 1455, BlockBytes: 28, Group1Block: 2, Group1BlockBytes: 46, Group2Block: 29, Group2BlockBytes: 47},
			{TotalBytes: 1033, BlockBytes: 30, Group1Block: 42, Group1BlockBytes: 24, Group2Block: 1, Group2BlockBytes: 25},
			{TotalBytes: 793, BlockBytes: 30, Group1Block: 23, Group1BlockBytes: 15, Group2Block: 28, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 1955, BlockBytes: 30, Group1Block: 17, Group1BlockBytes: 115, Group2Block: 0, Group2BlockBytes: 0},
			{TotalBytes: 1541, BlockBytes: 28, Group1Block: 10, Group1BlockBytes: 46, Group2Block: 23, Group2BlockBytes: 47},
			{TotalBytes: 1115, BlockBytes: 30, Group1Block: 10, Group1BlockBytes: 24, Group2Block: 35, Group2BlockBytes: 25},
			{TotalBytes: 845, BlockBytes: 30, Group1Block: 19, Group1BlockBytes: 15, Group2Block: 35, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 2071, BlockBytes: 30, Group1Block: 17, Group1BlockBytes: 115, Group2Block: 1, Group2BlockBytes: 116},
			{TotalBytes: 1631, BlockBytes: 28, Group1Block: 14, Group1BlockBytes: 46, Group2Block: 21, Group2BlockBytes: 47},
			{TotalBytes: 1171, BlockBytes: 30, Group1Block: 29, Group1BlockBytes: 24, Group2Block: 19, Group2BlockBytes: 25},
			{TotalBytes: 901, BlockBytes: 30, Group1Block: 11, Group1BlockBytes: 15, Group2Block: 46, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 2191, BlockBytes: 30, Group1Block: 13, Group1BlockBytes: 115, Group2Block: 6, Group2BlockBytes: 116},
			{TotalBytes: 1725, BlockBytes: 28, Group1Block: 14, Group1BlockBytes: 46, Group2Block: 23, Group2BlockBytes: 47},
			{TotalBytes: 1231, BlockBytes: 30, Group1Block: 44, Group1BlockBytes: 24, Group2Block: 7, Group2BlockBytes: 25},
			{TotalBytes: 961, BlockBytes: 30, Group1Block: 59, Group1BlockBytes: 16, Group2Block: 1, Group2BlockBytes: 17},
		},
		{
			{TotalBytes: 2306, BlockBytes: 30, Group1Block: 12, Group1BlockBytes: 121, Group2Block: 7, Group2BlockBytes: 122},
			{TotalBytes: 1812, BlockBytes: 28, Group1Block: 12, Group1BlockBytes: 47, Group2Block: 26, Group2BlockBytes: 48},
			{TotalBytes: 1286, BlockBytes: 30, Group1Block: 39, Group1BlockBytes: 24, Group2Block: 14, Group2BlockBytes: 25},
			{TotalBytes: 986, BlockBytes: 30, Group1Block: 22, Group1BlockBytes: 15, Group2Block: 41, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 2434, BlockBytes: 30, Group1Block: 6, Group1BlockBytes: 121, Group2Block: 14, Group2BlockBytes: 122},
			{TotalBytes: 1914, BlockBytes: 28, Group1Block: 6, Group1BlockBytes: 47, Group2Block: 34, Group2BlockBytes: 48},
			{TotalBytes: 1354, BlockBytes: 30, Group1Block: 46, Group1BlockBytes: 24, Group2Block: 10, Group2BlockBytes: 25},
			{TotalBytes: 1054, BlockBytes: 30, Group1Block: 2, Group1BlockBytes: 15, Group2Block: 64, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 2566, BlockBytes: 30, Group1Block: 17, Group1BlockBytes: 122, Group2Block: 4, Group2BlockBytes: 123},
			{TotalBytes: 1992, BlockBytes: 28, Group1Block: 29, Group1BlockBytes: 46, Group2Block: 14, Group2BlockBytes: 47},
			{TotalBytes: 1426, BlockBytes: 30, Group1Block: 49, Group1BlockBytes: 24, Group2Block: 10, Group2BlockBytes: 25},
			{TotalBytes: 1096, BlockBytes: 30, Group1Block: 24, Group1BlockBytes: 15, Group2Block: 46, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 2702, BlockBytes: 30, Group1Block: 4, Group1BlockBytes: 122, Group2Block: 18, Group2BlockBytes: 123},
			{TotalBytes: 2102, BlockBytes: 28, Group1Block: 13, Group1BlockBytes: 46, Group2Block: 32, Group2BlockBytes: 47},
			{TotalBytes: 1502, BlockBytes: 30, Group1Block: 48, Group1BlockBytes: 24, Group2Block: 14, Group2BlockBytes: 25},
			{TotalBytes: 1142, BlockBytes: 30, Group1Block: 42, Group1BlockBytes: 15, Group2Block: 32, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 2812, BlockBytes: 30, Group1Block: 20, Group1BlockBytes: 117, Group2Block: 4, Group2BlockBytes: 118},
			{TotalBytes: 2216, BlockBytes: 28, Group1Block: 40, Group1BlockBytes: 47, Group2Block: 7, Group2BlockBytes: 48},
			{TotalBytes: 1582, BlockBytes: 30, Group1Block: 43, Group1BlockBytes: 24, Group2Block: 22, Group2BlockBytes: 25},
			{TotalBytes: 1222, BlockBytes: 30, Group1Block: 10, Group1BlockBytes: 15, Group2Block: 67, Group2BlockBytes: 16},
		},
		{
			{TotalBytes: 2956, BlockBytes: 30, Group1Block: 19, Group1BlockBytes: 118, Group2Block: 6, Group2BlockBytes: 119},
			{TotalBytes: 2334, BlockBytes: 28, Group1Block: 18, Group1BlockBytes: 47, Group2Block: 31, Group2BlockBytes: 48},
			{TotalBytes: 1666, BlockBytes: 30, Group1Block: 34, Group1BlockBytes: 24, Group2Block: 34, Group2BlockBytes: 25},
			{TotalBytes: 1276, BlockBytes: 30, Group1Block: 20, Group1BlockBytes: 15, Group2Block: 61, Group2BlockBytes: 16},
		},
	}
)

type ecTable struct {
	TotalBytes       int // Total Number of Data Codewords for this Version and EC Level
	BlockBytes       int // EC Codewords Per Block
	Group1Block      int // Number of Blocks in Group 1
	Group1BlockBytes int // Number of Data Codewords in Each of Group 1's Blocks
	Group2Block      int // Number of Blocks in Group 2
	Group2BlockBytes int // Number of Data Codewords in Each of Group 2's Blocks
}

type version int

// 判断编码版本
func analysisVersion(level level, mode mode, str string) (version, error) {
	for i, a := range versionCapacity[level][mode] {
		if len(str) <= a {
			return version(i), nil
		}
	}
	return 0, fmt.Errorf("string length <%d> too lager", len(str))
}
