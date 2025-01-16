package nodes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/RulezKT/findf"
)

const RAD_TO_DEG = 5.7295779513082320877e1

const SEC_FILE = "nodes.sec"
const LNG_FILE = "nodes.lng"

type Nodes struct {
	secArr []int64
	lngArr []float64
	North  float64
	South  float64
}

func (n *Nodes) Load(folder string) {

	const FILE_LENGTH = 5397

	dir := findf.Dir(folder)

	secFile := findf.File(dir, SEC_FILE)
	f, err := os.ReadFile(secFile)
	if err != nil {
		log.Fatal(err)
	}

	r := bytes.NewReader(f)
	n.secArr = make([]int64, FILE_LENGTH)
	err = binary.Read(r, binary.LittleEndian, &n.secArr)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	lngFile := findf.File(dir, LNG_FILE)
	f, err = os.ReadFile(lngFile)
	if err != nil {
		log.Fatal(err)
	}

	r = bytes.NewReader(f)
	n.lngArr = make([]float64, FILE_LENGTH)
	err = binary.Read(r, binary.LittleEndian, &n.lngArr)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}

}

// считаем Лунные Узлы методом интерполяции
// V4 2024
// return in degrees
func (n *Nodes) Calc(dateInSeconds int64) {

	var startIndex int
	var nodeToFind float64

	for i, v := range n.secArr {
		if v > dateInSeconds {
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
	firstHalf := math.Abs(float64(dateInSeconds - startSecond))
	secondHalf := math.Abs(float64(endSecond - dateInSeconds))
	// считаем от 0 узла
	if firstHalf <= secondHalf {
		nodeToFind = startLng - nodeSpeed*firstHalf
		// считаем от узла +1
	} else {
		nodeToFind = endLng + nodeSpeed*secondHalf
	}
	// if x%2==0 if even then true else false
	// all even  indexes are north  all odd are south
	if startIndex%2 == 0 {
		n.North = nodeToFind
		n.South = n.North + 180
	} else {
		n.South = nodeToFind
		n.North = n.South + 180
	}

	if n.South > 360 {
		n.South -= 360
	}

	if n.North > 360 {
		n.North -= 360
	}

}
