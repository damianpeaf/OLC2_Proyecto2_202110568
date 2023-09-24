//struct con atributo sin valor por defecto
// y con un atributo con valor por defecto
struct Persona{
    var Nombre: String
    var edad = 0
}
// struct con funciones
struct Avion {
    var pasajeros = 0
    var velocidad = 100
    var nombre: String
    var piloto: Persona
    // metodo dentro de struct
    mutating func frenar(){
    print("Frenando")
    //al ser mutable sí afecta al struct
    self.velocidad = 0
    }
    // funcion inmutable
    func mostrarVelocidad(){
    print("Velocidad",self.velocidad)
    }
}

// creación de una instancia
var avioneta = Avion( nombre: "78496", piloto: Persona(Nombre: "Joel",edad: 43 ) )

// acceso a un atributo
print(avioneta.pasajeros)

// modificion de un atributo
avioneta.pasajeros = 5

print("Pasajeros:", avioneta.pasajeros)

// llamada de la funcion
avioneta.mostrarVelocidad()

// copia de structs por valor
var avioneta2 = avioneta
avioneta2.pasajeros = 0
//imprime: avioneta.pasajeros: 5
print("avioneta.pasajeros:",avioneta.pasajeros)
//avioneta2.pasajeros 0
print("avioneta2.pasajeros:",avioneta2.pasajeros)
print("avioneta.piloto.Nombre:",avioneta2.piloto.Nombre )
struct Fruta{
let nombre: String = "pera"
var precio: Int
}
// solo se puede definir precio en el constructor
// si se llega a definir nombre será un error
var pera = Fruta(precio: 10)
struct Verdura{
    let nombre: String
var precio: Int
}
// nombre se puede definir
//al no tener valor por defecto
var brocoli = Verdura(nombre:"brocoli", precio: 5)
struct Person {
var name: String
var age: Int
}
var personas = [Persona]()
// se agregan valores al arra
personas.append(Persona(Nombre: "Celeste", edad: 23))
personas.append(Persona(Nombre: "Roel", edad: 32))
personas.append(Persona(Nombre: "Flor", edad: 17))
// copia por valor
var persona1 = personas[0]
persona1.Nombre = "Nancy"
print(persona1.Nombre) //imprime Nancy
print(personas[0].Nombre) //imprime Celeste
// se modifica un array
personas[1].edad = 26
// otras formas permitidas
struct Distro {
var Nombre: String
var Version: String
}
var Distros = [
Distro(Nombre: "Ubuntu", Version: "22.04"),
Distro(Nombre: "Fedora", Version: "38"),
Distro(Nombre: "OpenSUSE", Version: "Leap 15")
]
//
print(Distros[0].Nombre) // Imprime Ubuntu
print(Distros[1].Version) // Imprime 13
// for con accesxo a los structs
for distro in Distros {
    print(distro.Nombre)
}
/* salida:
Ubuntu
Fedora
OpenSUSE
*/
// función que devuelve structs
func crearVerdura( precioV: Int, nombreV: String ) -> Verdura {
return Verdura( nombre: nombreV, precio: precioV )
}
// creación de struct por llmada de función
var verdura :Verdura = crearVerdura( precioV: 10, nombreV: "Apio")
// error por referencia indirecta
//este tipo de definciones es inválida
//Coordenada depende de Ubicacion
struct Coordenada{
var ubicacion: Ubicacion
var valorX: Int
var valorY: Int
}
// Ubicación depende de Coordenada
struct Ubicacion{
var nombre: String
var coordenada: Coordenada
}
// error por auto referencia
struct Nodo {
//no es posible que un struct tenga atributos de
//su mismo tipo
var siguiente: Nodo
}



struct Carro {
    let llantas: Int
}

var carro1 = Carro(llantas: 4)

func cambiarLlantas(_ car : inout Carro){
    car.llantas = 5
}

cambiarLlantas(&carro1)