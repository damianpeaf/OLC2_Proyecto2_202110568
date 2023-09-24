var w = Int("10")
var x =  Int("10.001")
var x1 = Int(10.999999)
var y = Int("Q10.00")

var res = x + w + x1 // 30
print(res)

var w = Float("10")
var x = Float("10.001")
var y = Float("Q10.00")

var res = x + w // 20.001
print(res)

print( String(10) + iota(3.5)) //imprime 103.5000
print( String( true )) //true
cadena = String(true) + "->" + String(3.504) //
print(cadena); // imprime true->3.50400000


var matrix : [[[Int]]] = [[[Int]]] (repeating: [[Int]] (repeating: [Int](repeating: 0, count:2), count:3), count:4)

for i in 0...3 {
    for j in 0...2 {
        for k in 0...1 {
            print(matrix[i][j][k])
        }
    }
}


// 

func imprimirArray (_ array: [Int] ) {
    for i in 0...array.count - 1 {
        print(array[i])
    }
}

func duplicarA (_ array: inout [Int] ) {
    var i = 0
    while (i < array.count ) {
        array[i] += array[i]
        i += 1
    }
}

func duplicarB (_ array: inout [Int] ) {
    duplicarA(&array)
}

func duplicarA2 (_ array: [Int] ) {
    for i in 0...array.count - 1 {
        array[i] += array[i]
    }
}


var array = [1,2,3,4,5,6]

print("Array original")
imprimirArray(array)

print("Array sin puntero")
duplicarA2(array)
imprimirArray(array)

print("Array con puntero")
duplicarA(&array)
imprimirArray(array)

print("Array con puntero doble")
duplicarB(&array)
imprimirArray(array)