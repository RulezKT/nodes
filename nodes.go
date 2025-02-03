package nodes

import (
	"math"
	"path/filepath"

	"github.com/RulezKT/floatsfile"
)

const RAD_TO_DEG = 5.7295779513082320877e1

const SEC_FILE = "nodes.sec"
const LNG_FILE = "nodes.lng"
const FILE_LENGTH = 5397

type Nodes struct {
	secArr []int64
	lngArr []float64
}

func (n *Nodes) Load(dir string) {

	n.secArr = floatsfile.LoadIntBinary(filepath.Join(dir, SEC_FILE), FILE_LENGTH)

	n.lngArr = floatsfile.LoadBinary(filepath.Join(dir, LNG_FILE), FILE_LENGTH)
}

// считаем Лунные Узлы методом интерполяции
// V4 2024
// return in degrees
func (n *Nodes) Calc(dateInSeconds float64) (float64, float64) {

	var startIndex int
	var nodeToFind float64

	for i, v := range n.secArr {
		if v > int64(dateInSeconds) {
			// fmt.Println("index =", i-1, "value = ", v)
			startIndex = i - 1
			break
		}
	}

	startSecond := n.secArr[startIndex]
	endSecond := n.secArr[startIndex+1]

	// находим начальную точку Узла, который считаем
	startLng := n.lngArr[startIndex]

	// находим финальную точку Узла, который считаем
	// Для этого берем позицию противоположного узла  через пол-месяца 27.2122/2 = 13.6061 дня.
	// и добавляем PI, так как узлы всегда находятся точно друг напротив друга
	endLng := n.lngArr[startIndex+1]
	// fmt.Println("node_clean_polar_end = ", node_clean_polar_end)
	endLng += 180
	// fmt.Println("node_clean_polar_end = ", node_clean_polar_end)
	if endLng > 360 {
		endLng -= 360
	}

	// fmt.Println("node_clean_polar_end = ", node_clean_polar_end)

	absDiff := math.Abs(endLng - startLng)
	if (absDiff) > 180+90 {
		if endLng > startLng {
			absDiff = 360 - endLng + startLng
		} else {
			absDiff = 360 - startLng + endLng
		}
	}
	// находим скорость передвижения узла за 1 секунду
	// для этого находим сколько прошел узел градусов за время прохода луны от одного узла до другого
	// примерно (27.2122/2 = 13.6061 дня.)
	nodeSpeed := absDiff / math.Abs(float64(endSecond-startSecond))

	// Проверка к какому из узлов ближе искомый узел и отсчитываем от него
	firstHalf := math.Abs(dateInSeconds - float64(startSecond))
	secondHalf := math.Abs(float64(endSecond) - dateInSeconds)
	// считаем от 0 узла
	if firstHalf <= secondHalf {
		nodeToFind = startLng - nodeSpeed*firstHalf
		// считаем от узла +1
	} else {
		nodeToFind = endLng + nodeSpeed*secondHalf
	}

	var northn, southn float64
	// if x%2==0 if even then true else false
	// all even  indexes are north  all odd are south
	if startIndex%2 == 0 {
		northn = nodeToFind
		southn = northn + 180
	} else {
		southn = nodeToFind
		northn = southn + 180
	}

	if southn > 360 {
		southn -= 360
	}

	if northn > 360 {
		northn -= 360
	}

	return northn, southn

}
