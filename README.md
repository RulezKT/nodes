    // its just an array of integers ,
    // each integer is the time of moon crossing one of the nodes
    // [0] index is north node
    // [1] index is south node
    // all even indexes are north
    // all odd are south

    // reworking nodes to have second value - moon longitude
    // [0][0] time
    // [0][1] longitude in Radians

    // fmt.Println("nodesCoords = ", len(nodesSec), len(nodesLng))

    // looking for the first index , that is greater then dateInSeconds
    // end_i = np.argmax(nodesCoords[:, 0] >= dateInSeconds)
