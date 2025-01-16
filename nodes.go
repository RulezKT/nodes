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
	SecArr []int64
	LngArr []float64
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
	n.SecArr = make([]int64, FILE_LENGTH)
	err = binary.Read(r, binary.LittleEndian, &n.SecArr)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	lngFile := findf.File(dir, LNG_FILE)
	f, err = os.ReadFile(lngFile)
	if err != nil {
		log.Fatal(err)
	}

	r = bytes.NewReader(f)
	n.LngArr = make([]float64, FILE_LENGTH)
	err = binary.Read(r, binary.LittleEndian, &n.LngArr)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}

}

// считаем Лунные Узлы методом интерполяции
// V4 2024
// return in degrees
func (n *Nodes) Calc(dateInSeconds int64) {

	var start_i int
	var node_to_find float64

	for i, v := range n.SecArr {
		if v > dateInSeconds {
			// fmt.Println("index =", i-1, "value = ", v)
			start_i = i - 1
			break
		}
	}

	start_second := n.SecArr[start_i]
	end_second := n.SecArr[start_i+1]

	// находим начальную точку Узла, который считаем
	node_clean_polar_start := n.LngArr[start_i]

	// находим финальную точку Узла, который считаем
	// Для этого берем позицию противоположного узла  через пол-месяца 27.2122/2 = 13.6061 дня.
	// и добавляем PI, так как узлы всегда находятся точно друг напротив друга
	node_clean_polar_end := n.LngArr[start_i+1]
	// fmt.Println("node_clean_polar_end = ", node_clean_polar_end)
	node_clean_polar_end += 180
	// fmt.Println("node_clean_polar_end = ", node_clean_polar_end)
	if node_clean_polar_end > 360 {
		node_clean_polar_end -= 360
	}

	// fmt.Println("node_clean_polar_end = ", node_clean_polar_end)

	abs_diff := math.Abs(node_clean_polar_end - node_clean_polar_start)
	if (abs_diff) > 180+90 {
		if node_clean_polar_end > node_clean_polar_start {
			abs_diff = 360 - node_clean_polar_end + node_clean_polar_start
		} else {
			abs_diff = 360 - node_clean_polar_start + node_clean_polar_end
		}
	}
	// находим скорость передвижения узла за 1 секунду
	// для этого находим сколько прошел узел градусов за время прохода луны от одного узла до другого
	// примерно (27.2122/2 = 13.6061 дня.)
	speed_of_node := abs_diff / math.Abs(float64(end_second-start_second))

	// Проверка к какому из узлов ближе искомый узел и отсчитываем от него
	first_halve_sec := math.Abs(float64(dateInSeconds - start_second))
	second_halve_sec := math.Abs(float64(end_second - dateInSeconds))
	// считаем от 0 узла
	if first_halve_sec <= second_halve_sec {
		node_to_find = node_clean_polar_start - speed_of_node*first_halve_sec
		// считаем от узла +1
	} else {
		node_to_find = node_clean_polar_end + speed_of_node*second_halve_sec
	}
	// if x%2==0 if even then true else false
	// all even  indexes are north  all odd are south
	if start_i%2 == 0 {
		n.North = node_to_find
		n.South = n.North + 180
	} else {
		n.South = node_to_find
		n.North = n.South + 180
	}

	if n.South > 360 {
		n.South -= 360
	}

	if n.North > 360 {
		n.North -= 360
	}

}
