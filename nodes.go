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

func Load(folder string) ([]int64, []float64) {

	const FILE_LENGTH = 5397

	dir := findf.Dir(folder)
	secFile := findf.File(dir, SEC_FILE)

	f, err := os.ReadFile(secFile)
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	r := bytes.NewReader(f)
	secArr := make([]int64, FILE_LENGTH)

	err = binary.Read(r, binary.LittleEndian, &secArr)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	lngFile := findf.File(dir, LNG_FILE)

	f, err = os.ReadFile(lngFile)
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	r = bytes.NewReader(f)
	lngArr := make([]float64, FILE_LENGTH)

	err = binary.Read(r, binary.LittleEndian, &lngArr)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	// r := bytes.NewReader(f)
	// floatArr := make([]float64, floats64Num)

	// err = binary.Read(r, binary.LittleEndian, &floatArr)
	// if err != nil {
	// 	fmt.Println("binary.Read failed:", err)
	// }

	// var secArr []int64
	// var lngArr []float64
	// for i := 0; i < len(floatArr); i++ {
	// 	if i%2 == 0 {
	// 		secArr = append(secArr, int64(floatArr[i]))
	// 	} else {
	// 		lngArr = append(lngArr, floatArr[i])
	// 	}
	// }

	// // write to file
	// file, err := os.Create("../files/nodes.sec")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	// err = binary.Write(file, binary.LittleEndian, &secArr)
	// if err != nil {
	// 	fmt.Println("binary.Write failed:", err)
	// }

	// for i, v := range lngArr {
	// 	lngArr[i] = v * RAD_TO_DEG
	// 	if lngArr[i] > 360 {
	// 		lngArr[i] -= 360
	// 	}
	// }

	// // write to file
	// file, err := os.Create("../files/nodesDEG.lng")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	// err = binary.Write(file, binary.LittleEndian, &lngArr)
	// if err != nil {
	// 	fmt.Println("binary.Write failed:", err)
	// }

	return secArr, lngArr
}

// считаем Лунные Узлы методом интерполяции
// V4 2024
// return in degrees
func (n Nodes) Calc(dateInSeconds int64, nodesSec []int64, nodesLng []float64) (float64, float64) {

	var start_i int

	for i, v := range nodesSec {
		if v > dateInSeconds {
			// fmt.Println("index =", i-1, "value = ", arr[i-1])
			start_i = i - 1
			break
		}
	}

	var north_node float64
	var south_node float64
	var node_to_find float64

	start_second := nodesSec[start_i]
	end_second := nodesSec[start_i+1]

	// находим начальную точку Узла, который считаем
	node_clean_polar_start := nodesLng[start_i]

	// находим финальную точку Узла, который считаем
	// Для этого берем позицию противоположного узла  через пол-месяца 27.2122/2 = 13.6061 дня.
	// и добавляем PI, так как узлы всегда находятся точно друг напротив друга
	node_clean_polar_end := nodesLng[start_i+1]
	// fmt.Println("node_clean_polar_end = ", node_clean_polar_end)
	node_clean_polar_end += 180
	// fmt.Println("node_clean_polar_end = ", node_clean_polar_end)
	node_clean_polar_end = Convert_to_0_360_DEG(node_clean_polar_end)

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
		north_node = node_to_find
		south_node = north_node + 180
	} else {
		south_node = node_to_find
		north_node = south_node + 180
	}

	south_node = Convert_to_0_360_DEG(south_node)
	north_node = Convert_to_0_360_DEG(north_node)

	return north_node, south_node
}

// убирает минус или значения больше 360
// если угол должен быть от 0 до 360
// все в Градусах
func Convert_to_0_360_DEG(longitude float64) float64 {

	coeff := int(math.Abs(longitude / 360))
	// fmt.Printf("coeff = %f\n", coeff)

	if longitude < 0 {
		return longitude + float64(coeff*360+360)
	} else {
		return longitude - float64(coeff*360)
	}

}
