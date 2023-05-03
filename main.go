package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type MBR struct {
	mbr_tamano         []byte
	mbr_fecha_creacion []byte
	mbr_dsk_signature  []byte
	mbr_dsk_fit        []byte
	mbr_partition_1    Partition
	mbr_partition_2    Partition
	mbr_partition_3    Partition
	mbr_partition_4    Partition
}

type Partition struct {
	part_status []byte
	part_type   []byte
	part_fit    []byte
	part_start  []byte
	part_size   []byte
	part_name   []byte
}

type EBR struct {
	part_status []byte
	part_fit    []byte
	part_start  []byte
	part_size   []byte
	part_next   []byte
	part_name   []byte
}

type SuperBloque struct {
	s_filesystem_type   []byte
	s_inodes_count      []byte
	s_blocks_count      []byte
	s_free_blocks_count []byte
	s_free_inodes_count []byte
	s_mtime             []byte
	s_mnt_count         []byte
	s_magic             []byte
	s_inode_size        []byte
	s_block_size        []byte
	s_first_ino         []byte
	s_first_blo         []byte
	s_bm_inode_start    []byte
	s_bm_block_start    []byte
	s_inode_start       []byte
	s_block_start       []byte
}

type TablaInodo struct {
	i_uid   []byte
	i_gid   []byte
	i_size  []byte
	i_atime []byte
	i_ctime []byte
	i_mtime []byte
	i_block []byte
	i_type  []byte
	i_perm  []byte
}

type BloqueCarpeta struct {
	b_content [4]content
}

type content struct {
	b_name  []byte
	b_inodo []byte
}

type BloqueArchivos struct {
	b_content []byte
}

type MountPart struct {
	Active     bool
	ID         string
	Path       string
	part_name  []byte
	part_size  []byte
	part_start []byte
}

type Session struct {
	User   string
	IDU    int
	Grupo  string
	IDG    int
	Active MountPart
}

var ActivePart [10]MountPart
var Sesion Session
var Console string

func main() {
	Inicializacion()
	AnalizarCodigo(`mkdisk >size=50 >unit=m >path=/home/dabs/201807100/Disco1.dk
mkdisk >size=50 >unit=m >path=/home/dabs/201807100/Disco2.dk
rmdisk >path=/home/dabs/201807100/Disco2.dk

fdisk >size=5 >path=/home/dabs/201807100/Disco1.dk >unit=m >name=Particion1 >fit=ff
fdisk >S=1024 >path=/home/dabs/201807100/Disco1.dk >unit=k >name=Particion2
fdisk >size=1024 >path=/home/dabs/201807100/Disco1.dk >unit=k >name=Particion2
fdisk >size=10 >unit=m >path=/home/dabs/201807100/Disco1.dk >name=Particion3
fdisk >size=25 >path=/home/dabs/201807100/Disco1.dk >name=Particion4 >fit=wf >unit=m 
fdisk >size=25 >path=/home/dabs/201807100/Disco1.dk >name=Particion4 >fit=wf >unit=m
#seedisk >path=/home/dabs/201807100/Disco1.dk

mkdisk >size=25 >fit=bf >unit=m >path="/home/dabs/201807100/primer semestre/Disco2.dk"
fdisk >size=500 >unit=k >path="/home/dabs/201807100/primer semestre/Disco2.dk" >name=Particion1 >fit=ff
fdisk >size=1024 >path="/home/dabs/201807100/primer semestre/Disco2.dk" >unit=k >name=Particion2
fdisk >size=10 >unit=m >path="/home/dabs/201807100/primer semestre/Disco2.dk" >name=Particion3
fdisk >unit=k >size=4096 >path="/home/dabs/201807100/primer semestre/Disco2.dk" >type=E >name=Particion4 >fit=wf
#seedisk >path="/home/dabs/201807100/primer semestre/Disco2.dk"

mkdisk >unit=k >size=75 >path="/home/dabs/201807100/primer semestre/entrada1/Disco3.dk"
fdisk >size=5000 >path="/home/dabs/201807100/primer semestre/entrada1/Disco3.dk" >name=Particion1 >unit=b
fdisk >size=30 >path="/home/dabs/201807100/primer semestre/entrada1/Disco3.dk" >unit=m >type=E >fit=bf >name=Particion2
fdisk >size=5 >type=L >unit=m >path="/home/dabs/201807100/primer semestre/entrada1/Disco3.dk" >name=Particion3
fdisk >type=L >unit=k >size=4096 >path="/home/dabs/201807100/primer semestre/entrada1/Disco3.dk" >name=Particion4
fdisk >size=3 >path="/home/dabs/201807100/primer semestre/entrada1/Disco3.dk" >name=Particion5 >unit=m
seedisk >path="/home/dabs/201807100/primer semestre/entrada1/Disco3.dk"
mount >path=/home/dabs/201807100/Disco1.dk >name=Particion1
mount >path="/home/dabs/201807100/primer semestre/Disco2.dk" >name=Particion2
mount >path="/home/dabs/201807100/primer semestre/entrada1/Disco3.dk" >name=Particion3
#seemounts

mkfs >type=full >id=001a
login >pwd=123 >user=root >id=001a

mkgrp >name=prueba1
mkgrp >name=prueba2
mkgrp >name=prueba3
rmgrp >name=prueba3

mkusr >user="user1" >grp=prueba1 >pwf=user1
mkusr >user="user2" >grp=prueba2 >pwf=user2
mkusr >user="user3" >grp=prueba1 >pwf=user3
mkusr >user="user4" >grp=prueba2 >pwf=user4
rmusr >user=user4

mkfile >size=10 >R >path=/home/archivos/dabs/fase2/docs/a.txt
mkfile >size=0 >path=/home/archivos/dabs/fase2/docs/b.txt
mkfile >size=20 >path=/home/archivos/dabs/fase2/docs/c.txt
mkfile >path="/home/dabs/pruebas/hola.txt" >r >cont="/home/dabs/Documentos/GitHub/-MIA-Proyecto1_201807100/archivo2.txt"

seeinfo >id=001a

logout`)
	fmt.Println("Consola: ")
	fmt.Println(Console[:len(Console)-1])
}

func AnalizarCodigo(Entrada string) {
	Lineas := strings.Split(Entrada, "\n")
	i := 0
	for i < len(Lineas) && LeerComando(Lineas[i]+" ") {
		i++
	}
}

func LeerComando(Linea string) bool {
	Error := false
	Aux := Valor(&Linea)
	var err error
	if Aux != "" && Aux[0] != '#' {
		if strings.ToLower(Aux) == "mkdisk" {
			Size := 0
			Path := ""
			Fit := "FF"
			Unit := "M"
			BSize := false
			BPath := false
			for Linea != "" && !Error {
				Aux = Valor(&Linea)
				if strings.ToLower(Aux) == ">size" {
					Size, err = strconv.Atoi(Valor(&Linea))
					if err != nil {
						Error = true
					}
					if Size <= 0 {
						Error = true
						Console += "Error, El tamaño del disco debe ser mayor a 0\n"
					}
					BSize = true
				} else if strings.ToLower(Aux) == ">path" {
					Path = Valor(&Linea)
					BPath = true
				} else if strings.ToLower(Aux) == ">fit" {
					Fit = Valor(&Linea)
					if !(strings.ToLower(Fit) != "bf" || strings.ToLower(Fit) != "ff" || strings.ToLower(Fit) != "wf") {
						Error = true
						Console += "Error, Tipo de fit invalido o no soportado\n"
					}
				} else if strings.ToLower(Aux) == ">unit" {
					Unit = Valor(&Linea)
					if !(strings.ToLower(Unit) != "k" || strings.ToLower(Unit) != "m") {
						Error = true
						Console += "Error, Unidad invalido o no soportado\n"
					}
				} else {
					Error = true
					Console += "Error, Parametro desconocido\n"
				}
			}
			if BSize && BPath && !Error {
				CrearDisco(Size, Path, Fit, Unit)
			} else if !Error {
				Console += "Error, Faltan parametros\n"
			}
		} else if strings.ToLower(Aux) == "rmdisk" {
			Path := ""
			BPath := false
			for Linea != "" && !Error {
				Aux = Valor(&Linea)
				if strings.ToLower(Aux) == ">path" {
					Path = Valor(&Linea)
					BPath = true
				} else {
					Error = true
					Console += "Error, Parametro desconocido\n"
				}
			}
			if BPath && !Error {
				EliminarDisco(Path)
			} else if !Error {
				Console += "Error, Faltan parametros\n"
			}
		} else if strings.ToLower(Aux) == "fdisk" {
			Size := 0
			Path := ""
			Fit := "FF"
			Unit := "M"
			Type := "P"
			Name := ""
			BSize := false
			BPath := false
			BName := false
			for Linea != "" && !Error {
				Aux = Valor(&Linea)
				if strings.ToLower(Aux) == ">size" {
					Size, err = strconv.Atoi(Valor(&Linea))
					if err != nil {
						Error = true
					}
					if Size <= 0 {
						Error = true
						Console += "Error, El tamaño del disco debe ser mayor a 0\n"
					}
					BSize = true
				} else if strings.ToLower(Aux) == ">path" {
					Path = Valor(&Linea)
					BPath = true
				} else if strings.ToLower(Aux) == ">fit" {
					Fit = Valor(&Linea)
					if !(strings.ToLower(Fit) != "bf" || strings.ToLower(Fit) != "ff" || strings.ToLower(Fit) != "wf") {
						Error = true
						Console += "Error, Tipo de fit invalido o no soportado\n"
					}
				} else if strings.ToLower(Aux) == ">unit" {
					Unit = Valor(&Linea)
					if !(strings.ToLower(Unit) != "k" || strings.ToLower(Unit) != "m") {
						Error = true
						Console += "Error, Unidad invalido o no soportado\n"
					}
				} else if strings.ToLower(Aux) == ">type" {
					Type = Valor(&Linea)
					if !(strings.ToLower(Type) != "p" || strings.ToLower(Type) != "e" || strings.ToLower(Type) != "l") {
						Error = true
						Console += "Error, Unidad invalido o no soportado\n"
					}
				} else if strings.ToLower(Aux) == ">name" {
					Name = Valor(&Linea)
					if len(Name) > 16 {
						Error = true
						Console += "Error, El nombre debe ser como maximo de 16 caracteres"
					}
					BName = true
				} else {
					Error = true
					Console += "Error, Parametro desconocido\n"
				}
			}
			if BSize && BPath && BName && !Error {
				CrearPartición(Size, Path, Name, Unit, Type, Fit)
			} else if !Error {
				Console += "Error, Faltan parametros\n"
			}
		} else if strings.ToLower(Aux) == "mount" {
			Path := ""
			Name := ""
			BPath := false
			BName := false
			for Linea != "" && !Error {
				Aux = Valor(&Linea)
				if strings.ToLower(Aux) == ">path" {
					Path = Valor(&Linea)
					BPath = true
				} else if strings.ToLower(Aux) == ">name" {
					Name = Valor(&Linea)
					if len(Name) > 16 {
						Error = true
						Console += "Error, El nombre debe ser como maximo de 16 caracteres"
					}
					BName = true
				} else {
					Error = true
					Console += "Error, Parametro desconocido\n"
				}
			}
			if BPath && BName && !Error {
				MountParticion(Path, Name)
			} else if !Error {
				Console += "Error, Faltan parametros\n"
			}
		} else if strings.ToLower(Aux) == "mkfs" {
			ID := ""
			Type := ""
			BID := false
			for Linea != "" && !Error {
				Aux = Valor(&Linea)
				if strings.ToLower(Aux) == ">id" {
					ID = Valor(&Linea)
					BID = true
				} else if strings.ToLower(Aux) == ">type" {
					Type = Valor(&Linea)
					if strings.ToLower(Type) != "full" {
						Error = true
						Console += "Error, Opcion invalida o no soportada"
					}
				} else {
					Error = true
					Console += "Error, Parametro desconocido\n"
				}
			}
			if BID && !Error {
				FormatearParticion(ID)
			} else if !Error {
				Console += "Error, Faltan parametros\n"
			}
		} else if strings.ToLower(Aux) == "login" {
			if Sesion.User == "" {
				User := ""
				Pass := ""
				ID := ""
				BUser := false
				BPass := false
				BID := false
				for Linea != "" && !Error {
					Aux = Valor(&Linea)
					if strings.ToLower(Aux) == ">user" {
						User = Valor(&Linea)
						BUser = true
					} else if strings.ToLower(Aux) == ">pwd" {
						Pass = Valor(&Linea)
						BPass = true
					} else if strings.ToLower(Aux) == ">id" {
						ID = Valor(&Linea)
						BID = true
					} else {
						Error = true
						Console += "Error, Parametro desconocido\n"
					}
				}
				if BUser && BPass && BID && !Error {
					Login(User, Pass, ID)
				} else if !Error {
					Console += "Error, Faltan parametros\n"
				}
			} else {
				Console += "Error, ya esta logeado\n"
			}
		} else if strings.ToLower(Aux) == "logout" {
			if Sesion.User != "" {
				Sesion.User = ""
				Sesion.Grupo = ""
				Sesion.IDU = -1
				Sesion.IDG = -1
				Sesion.Active.Active = false
				Console += "Sesion Cerrada Exitosamente\n"
			} else {
				Console += "Error, usted no esta logeado\n"
			}
		} else if strings.ToLower(Aux) == "mkgrp" {
			if Sesion.User == "root      " {
				Name := ""
				BName := false
				for Linea != "" && !Error {
					Aux = Valor(&Linea)
					if strings.ToLower(Aux) == ">name" {
						Name = Valor(&Linea)
						if len(Name) < 1 {
							Console += "Error, no puede ingresar un nombre vacio\n"
							Error = true
						} else if len(Name) > 10 {
							Console += "Error, el nombre es demaciado largo\n"
							Error = true
						}
						BName = true
					} else {
						Error = true
						Console += "Error, Parametro desconocido\n"
					}
				}
				if BName && !Error {
					AgregarGrupo(Name)
				} else if !Error {
					Console += "Error, Faltan parametros\n"
				}
			} else if Sesion.User == "" {
				Console += "Error, No tiene sesion activa\n"
			} else {
				Console += "Error, Solo el usuario root puede agregar grupos\n"
			}
		} else if strings.ToLower(Aux) == "rmgrp" {
			if Sesion.User == "root      " {
				Name := ""
				BName := false
				for Linea != "" && !Error {
					Aux = Valor(&Linea)
					if strings.ToLower(Aux) == ">name" {
						Name = Valor(&Linea)
						if len(Name) < 1 {
							Console += "Error, no puede ingresar un nombre vacio\n"
							Error = true
						} else if len(Name) > 10 {
							Console += "Error, el nombre es demaciado largo\n"
							Error = true
						}
						BName = true
					} else {
						Error = true
						Console += "Error, Parametro desconocido\n"
					}
				}
				if BName && !Error {
					EliminarGrupo(Name)
				} else if !Error {
					Console += "Error, Faltan parametros\n"
				}
			} else if Sesion.User == "" {
				Console += "Error, No tiene sesion activa\n"
			} else {
				Console += "Error, Solo el usuario root puede agregar grupos\n"
			}
		} else if strings.ToLower(Aux) == "mkusr" {
			if Sesion.User == "root      " {
				User := ""
				Pwd := ""
				Grp := ""
				BUser := false
				BPwd := false
				BGrp := false
				for Linea != "" && !Error {
					Aux = Valor(&Linea)
					if strings.ToLower(Aux) == ">user" {
						User = Valor(&Linea)
						if len(User) < 1 {
							Console += "Error, no puede ingresar un usuario vacio\n"
							Error = true
						} else if len(User) > 10 {
							Console += "Error, el usuario es demaciado largo\n"
							Error = true
						}
						BUser = true
					} else if strings.ToLower(Aux) == ">pwf" {
						Pwd = Valor(&Linea)
						if len(Pwd) < 1 {
							Console += "Error, no puede ingresar una contraseña vacia\n"
							Error = true
						} else if len(Pwd) > 10 {
							Console += "Error, la contraseña es demaciado larga\n"
							Error = true
						}
						BPwd = true
					} else if strings.ToLower(Aux) == ">grp" {
						Grp = Valor(&Linea)
						if len(Grp) < 1 {
							Console += "Error, no puede ingresar un grupo vacio\n"
							Error = true
						} else if len(Grp) > 10 {
							Console += "Error, el nombre de grupo es demaciado largo\n"
							Error = true
						}
						BGrp = true
					} else {
						Error = true
						Console += "Error, Parametro desconocido\n"
					}
				}
				if BUser && BPwd && BGrp && !Error {
					AgregarUsuario(User, Pwd, Grp)
				} else if !Error {
					Console += "Error, Faltan parametros\n"
				}
			} else if Sesion.User == "" {
				Console += "Error, No tiene sesion activa\n"
			} else {
				Console += "Error, Solo el usuario root puede agregar grupos\n"
			}
		} else if strings.ToLower(Aux) == "rmusr" {
			if Sesion.User == "root      " {
				User := ""
				BUser := false
				for Linea != "" && !Error {
					Aux = Valor(&Linea)
					if strings.ToLower(Aux) == ">user" {
						User = Valor(&Linea)
						if len(User) < 1 {
							Console += "Error, no puede ingresar un nombre vacio\n"
							Error = true
						} else if len(User) > 10 {
							Console += "Error, el nombre es demaciado largo\n"
							Error = true
						}
						BUser = true
					} else {
						Error = true
						Console += "Error, Parametro desconocido\n"
					}
				}
				if BUser && !Error {
					EliminarUsuario(User)
				} else if !Error {
					Console += "Error, Faltan parametros\n"
				}
			} else if Sesion.User == "" {
				Console += "Error, No tiene sesion activa\n"
			} else {
				Console += "Error, Solo el usuario root puede agregar grupos\n"
			}
		} else if strings.ToLower(Aux) == "mkfile" {
			if Sesion.User != "" {
				Path := ""
				Cont := ""
				Size := 0
				BPath := false
				R := false
				for Linea != "" && !Error {
					Aux = Valor(&Linea)
					if strings.ToLower(Aux) == ">path" {
						Path = Valor(&Linea)
						BPath = true
					} else if strings.ToLower(Aux) == ">cont" {
						Cont = Valor(&Linea)
					} else if strings.ToLower(Aux) == ">size" {
						Size, err = strconv.Atoi(Valor(&Linea))
						if err != nil {
							Error = true
						}
						if Size < 0 {
							Error = true
							Console += "Error, El tamaño del disco debe ser mayor a 0\n"
						}
					} else if strings.ToLower(Aux) == ">r" {
						R = true
					} else {
						Error = true
						Console += "Error, Parametro desconocido\n"
					}
				}
				if BPath && !Error {
					MKFILE(Path, R, Size, Cont)
				} else if !Error {
					Console += "Error, Faltan parametros\n"
				}
			} else {
				Console += "Error, No tiene sesion activa\n"
			}
		} else if strings.ToLower(Aux) == "pause" {
			fmt.Println("Consola: ")
			fmt.Println(Console[:len(Console)-1])
			Console = ""
			fmt.Print("Ejecución en pausa, oprima enter para continuar ")
			fmt.Scanln()
		} else if strings.ToLower(Aux) == "seedisk" {
			Path := ""
			BPath := false
			for Linea != "" && !Error {
				Aux = Valor(&Linea)
				if strings.ToLower(Aux) == ">path" {
					Path = Valor(&Linea)
					BPath = true
				} else {
					Error = true
					Console += "Error, Parametro desconocido\n"
				}
			}
			if BPath && !Error {
				VerDisco(Path)
			} else if !Error {
				Console += "Error, Faltan parametros\n"
			}
		} else if strings.ToLower(Aux) == "seemounts" {
			VerMounts()
		} else if strings.ToLower(Aux) == "seeinfo" {
			ID := ""
			BID := false
			for Linea != "" && !Error {
				Aux = Valor(&Linea)
				if strings.ToLower(Aux) == ">id" {
					ID = Valor(&Linea)
					BID = true
				} else {
					Error = true
					Console += "Error, Parametro desconocido\n"
				}
			}
			if BID && !Error {
				VerInfo(ID)
			} else if !Error {
				Console += "Error, Faltan parametros\n"
			}
		} else {
			Console += "Error, comando desconocido\n"
		}
	}
	return true
}

//Funciones de comando

func CrearDisco(Size int, Path string, Fit string, Unit string) {
	var err error
	var NewMBR MBR
	InitMBR(&NewMBR)
	NewMBR.mbr_fecha_creacion, err = (time.Now()).MarshalBinary()
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		Console += "Error al optener la fecha\n"
	}
	binary.LittleEndian.PutUint32(NewMBR.mbr_dsk_signature, uint32(rand.Intn(101)))
	if strings.ToLower(Fit) == "bf" {
		NewMBR.mbr_dsk_fit[0] = 'B'
	} else if strings.ToLower(Fit) == "ff" {
		NewMBR.mbr_dsk_fit[0] = 'F'
	} else {
		NewMBR.mbr_dsk_fit[0] = 'W'
	}
	NewMBR.mbr_partition_1.part_status[0] = 'D'
	NewMBR.mbr_partition_1.part_type[0] = ' '
	NewMBR.mbr_partition_1.part_fit[0] = ' '
	binary.LittleEndian.PutUint32(NewMBR.mbr_partition_1.part_start, uint32(0))
	binary.LittleEndian.PutUint32(NewMBR.mbr_partition_1.part_size, uint32(0))
	copy(NewMBR.mbr_partition_1.part_name[:], []byte(""))
	NewMBR.mbr_partition_2.part_status[0] = 'D'
	NewMBR.mbr_partition_2.part_type[0] = ' '
	NewMBR.mbr_partition_2.part_fit[0] = ' '
	binary.LittleEndian.PutUint32(NewMBR.mbr_partition_2.part_start, uint32(0))
	binary.LittleEndian.PutUint32(NewMBR.mbr_partition_2.part_size, uint32(0))
	copy(NewMBR.mbr_partition_2.part_name[:], []byte(""))
	NewMBR.mbr_partition_3.part_status[0] = 'D'
	NewMBR.mbr_partition_3.part_type[0] = ' '
	NewMBR.mbr_partition_3.part_fit[0] = ' '
	binary.LittleEndian.PutUint32(NewMBR.mbr_partition_3.part_start, uint32(0))
	binary.LittleEndian.PutUint32(NewMBR.mbr_partition_3.part_size, uint32(0))
	copy(NewMBR.mbr_partition_3.part_name[:], []byte(""))
	NewMBR.mbr_partition_4.part_status[0] = 'D'
	NewMBR.mbr_partition_4.part_type[0] = ' '
	NewMBR.mbr_partition_4.part_fit[0] = ' '
	binary.LittleEndian.PutUint32(NewMBR.mbr_partition_4.part_start, uint32(0))
	binary.LittleEndian.PutUint32(NewMBR.mbr_partition_4.part_size, uint32(0))
	copy(NewMBR.mbr_partition_4.part_name[:], []byte(""))
	CrearPath(Path)
	archivo, err := os.OpenFile(Path, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println(err)
	}
	defer archivo.Close()
	var Data []byte
	if strings.ToLower(Unit) == "K" {
		binary.LittleEndian.PutUint32(NewMBR.mbr_tamano, uint32(Size*1024))
		Data = make([]byte, 1024)
	} else {
		binary.LittleEndian.PutUint32(NewMBR.mbr_tamano, uint32(Size*1024*1024))
		Data = make([]byte, 1024*1024)
	}
	for i := range Data {
		Data[i] = 0
	}
	for i := 0; i < Size; i++ {
		err = binary.Write(archivo, binary.LittleEndian, Data)
	}
	EscribirMBR(Path, NewMBR)
}

func EliminarDisco(Path string) {
	os.Remove(Path)
	Console += "Archivo Eliminado con exito\n"
}

func CrearPartición(Size int, Path string, Name string, Unit string, Type string, Fit string) {
	if strings.ToLower(Unit) == "k" {
		Size = Size * 1024
	} else if strings.ToLower(Unit) == "m" {
		Size = Size * 1024 * 1024
	}
	var MBRDsk MBR
	InitMBR(&MBRDsk)
	LeerMBR(Path, &MBRDsk)
	SizeMBR := int(unsafe.Sizeof(MBRDsk))
	if strings.ToLower(Type) == "p" || strings.ToLower(Type) == "e" {
		if strings.ToLower(Type) == "p" || (string(MBRDsk.mbr_partition_1.part_type) != "E" && string(MBRDsk.mbr_partition_2.part_type) != "E" && string(MBRDsk.mbr_partition_3.part_type) != "E" && string(MBRDsk.mbr_partition_4.part_type) != "E") {
			if string(MBRDsk.mbr_partition_4.part_status) == "D" {
				RegistrarEn := 5
				SizeDsk := int(binary.LittleEndian.Uint32(MBRDsk.mbr_tamano))
				Part1Start := int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_1.part_start))
				Part1Size := int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_1.part_size))
				Part2Start := int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_2.part_start))
				Part2Size := int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_2.part_size))
				Part3Start := int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_3.part_start))
				Part3Size := int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_3.part_size))
				TypeFit := ' '
				if strings.ToLower(Fit) == "bf" {
					TypeFit = 'B'
					if string(MBRDsk.mbr_partition_1.part_status) == "D" {
						if (SizeDsk - SizeMBR) >= Size {
							RegistrarEn = 1
						}
					} else {
						MinSize := SizeDsk
						if (Part1Start - SizeMBR) >= Size {
							MinSize = Part1Start - SizeMBR
							RegistrarEn = 1
						}
						if string(MBRDsk.mbr_partition_2.part_status) == "D" && SizeDsk-(Part1Start+Part1Size) >= Size && MinSize > SizeDsk-(Part1Start+Part1Size) {
							RegistrarEn = 2
						} else {
							if Part2Start-(Part1Start+Part1Size) > Size && MinSize > Part2Start-(Part1Start+Part1Size) {
								MinSize = Part2Start - (Part1Start + Part1Size)
								RegistrarEn = 2
							}
							if string(MBRDsk.mbr_partition_3.part_status) == "D" && SizeDsk-(Part2Start+Part2Size) >= Size && MinSize > SizeDsk-(Part2Start+Part2Size) {
								RegistrarEn = 3
							} else {
								if Part3Start-(Part2Start+Part2Size) >= Size && MinSize > Part3Start-(Part2Start+Part2Size) {
									MinSize = Part3Start - (Part2Start + Part2Size)
									RegistrarEn = 3
								}
								if SizeDsk-(Part3Start+Part3Size) >= Size && MinSize > SizeDsk-(Part3Start+Part3Size) {
									RegistrarEn = 4
								}
							}
						}
					}
				} else if strings.ToLower(Fit) == "ff" {
					TypeFit = 'F'
					if (string(MBRDsk.mbr_partition_1.part_status) == "D" && SizeDsk-SizeMBR >= Size) || Part1Start-SizeMBR >= Size {
						RegistrarEn = 1
					} else if (string(MBRDsk.mbr_partition_2.part_status) == "D" && SizeDsk-(Part1Start+Part1Size) >= Size) || Part2Start-(Part1Start+Part1Size) >= Size {
						RegistrarEn = 2
					} else if (string(MBRDsk.mbr_partition_3.part_status) == "D" && SizeDsk-(Part2Start+Part2Size) >= Size) || Part3Start-(Part2Start+Part2Size) >= Size {
						RegistrarEn = 3
					} else if SizeDsk-(Part3Start+Part3Size) >= Size {
						RegistrarEn = 4
					}
				} else {
					TypeFit = 'W'
					if string(MBRDsk.mbr_partition_1.part_status) == "D" {
						if (SizeDsk - SizeMBR) >= Size {
							RegistrarEn = 1
						}
					} else {
						MaxSize := 0
						if (Part1Start - SizeMBR) >= Size {
							MaxSize = Part1Start - SizeMBR
							RegistrarEn = 1
						}
						if string(MBRDsk.mbr_partition_2.part_status) == "D" && SizeDsk-(Part1Start+Part1Size) >= Size && MaxSize < SizeDsk-(Part1Start+Part1Size) {
							RegistrarEn = 2
						} else {
							if Part2Start-(Part1Start+Part1Size) > Size && MaxSize < Part2Start-(Part1Start+Part1Size) {
								MaxSize = Part2Start - (Part1Start + Part1Size)
								RegistrarEn = 2
							}
							if string(MBRDsk.mbr_partition_3.part_status) == "D" && SizeDsk-(Part2Start+Part2Size) >= Size && MaxSize < SizeDsk-(Part2Start+Part2Size) {
								RegistrarEn = 3
							} else {
								if Part3Start-(Part2Start+Part2Size) >= Size && MaxSize < Part3Start-(Part2Start+Part2Size) {
									MaxSize = Part3Start - (Part2Start + Part2Size)
									RegistrarEn = 3
								}
								if SizeDsk-(Part3Start+Part3Size) >= Size && MaxSize < SizeDsk-(Part3Start+Part3Size) {
									RegistrarEn = 4
								}
							}
						}
					}
				}
				Registrado := true
				switch RegistrarEn {
				case 1:
					MBRDsk.mbr_partition_4 = MBRDsk.mbr_partition_3.clone()
					MBRDsk.mbr_partition_3 = MBRDsk.mbr_partition_2.clone()
					MBRDsk.mbr_partition_2 = MBRDsk.mbr_partition_1.clone()
					MBRDsk.mbr_partition_1.part_status[0] = 'A'
					MBRDsk.mbr_partition_1.part_fit[0] = byte(TypeFit)
					binary.LittleEndian.PutUint32(MBRDsk.mbr_partition_1.part_start, uint32(SizeMBR))
					binary.LittleEndian.PutUint32(MBRDsk.mbr_partition_1.part_size, uint32(Size))
					copy(MBRDsk.mbr_partition_1.part_name[:], []byte(Name))
					if strings.ToLower(Type) == "p" {
						MBRDsk.mbr_partition_1.part_type[0] = byte('P')
					} else {
						var NewEBR EBR
						InitEBR(&NewEBR)
						if Size > int(unsafe.Sizeof(NewEBR)) {
							MBRDsk.mbr_partition_1.part_type[0] = byte('E')
							NewEBR.part_status[0] = 'D'
							NewEBR.part_fit[0] = ' '
							NewEBR.part_start = MBRDsk.mbr_partition_1.part_start
							binary.LittleEndian.PutUint32(NewEBR.part_size, uint32(0))
							binary.LittleEndian.PutUint32(NewEBR.part_next, uint32(0))
							copy(NewEBR.part_name, []byte(""))
							EscribirEBR(Path, NewEBR, SizeMBR)
						} else {
							Registrado = false
							Console += "Error, no hay espacio para crear un EBR en la particion Extendida\n"
						}
					}
				case 2:
					MBRDsk.mbr_partition_4 = MBRDsk.mbr_partition_3.clone()
					MBRDsk.mbr_partition_3 = MBRDsk.mbr_partition_2.clone()
					MBRDsk.mbr_partition_2.part_status[0] = 'A'
					MBRDsk.mbr_partition_2.part_fit[0] = byte(TypeFit)
					binary.LittleEndian.PutUint32(MBRDsk.mbr_partition_2.part_start, uint32(Part1Start+Part1Size))
					binary.LittleEndian.PutUint32(MBRDsk.mbr_partition_2.part_size, uint32(Size))
					copy(MBRDsk.mbr_partition_2.part_name[:], []byte(Name))
					if strings.ToLower(Type) == "p" {
						MBRDsk.mbr_partition_2.part_type[0] = byte('P')
					} else {
						var NewEBR EBR
						InitEBR(&NewEBR)
						if Size > int(unsafe.Sizeof(NewEBR)) {
							MBRDsk.mbr_partition_2.part_type[0] = byte('E')
							NewEBR.part_status[0] = 'D'
							NewEBR.part_fit[0] = ' '
							NewEBR.part_start = MBRDsk.mbr_partition_2.part_start
							binary.LittleEndian.PutUint32(NewEBR.part_size, uint32(0))
							binary.LittleEndian.PutUint32(NewEBR.part_next, uint32(0))
							copy(NewEBR.part_name, []byte(""))
							EscribirEBR(Path, NewEBR, Part1Start+Part1Size)
						} else {
							Registrado = false
							Console += "Error, no hay espacio para crear un EBR en la particion Extendida\n"
						}
					}
				case 3:
					MBRDsk.mbr_partition_4 = MBRDsk.mbr_partition_3.clone()
					MBRDsk.mbr_partition_3.part_status[0] = 'A'
					MBRDsk.mbr_partition_3.part_fit[0] = byte(TypeFit)
					binary.LittleEndian.PutUint32(MBRDsk.mbr_partition_3.part_start, uint32(Part2Start+Part2Size))
					binary.LittleEndian.PutUint32(MBRDsk.mbr_partition_3.part_size, uint32(Size))
					copy(MBRDsk.mbr_partition_3.part_name[:], []byte(Name))
					if strings.ToLower(Type) == "p" {
						MBRDsk.mbr_partition_3.part_type[0] = byte('P')
					} else {
						var NewEBR EBR
						InitEBR(&NewEBR)
						if Size > int(unsafe.Sizeof(NewEBR)) {
							MBRDsk.mbr_partition_3.part_type[0] = byte('E')
							NewEBR.part_status[0] = 'D'
							NewEBR.part_fit[0] = ' '
							NewEBR.part_start = MBRDsk.mbr_partition_3.part_start
							binary.LittleEndian.PutUint32(NewEBR.part_size, uint32(0))
							binary.LittleEndian.PutUint32(NewEBR.part_next, uint32(0))
							copy(NewEBR.part_name, []byte(""))
							EscribirEBR(Path, NewEBR, Part2Start+Part2Size)
						} else {
							Registrado = false
							Console += "Error, no hay espacio para crear un EBR en la particion Extendida\n"
						}
					}
				case 4:
					MBRDsk.mbr_partition_4.part_status[0] = 'A'
					MBRDsk.mbr_partition_4.part_fit[0] = byte(TypeFit)
					binary.LittleEndian.PutUint32(MBRDsk.mbr_partition_4.part_start, uint32(Part3Start+Part3Size))
					binary.LittleEndian.PutUint32(MBRDsk.mbr_partition_4.part_size, uint32(Size))
					copy(MBRDsk.mbr_partition_4.part_name[:], []byte(Name))
					if strings.ToLower(Type) == "p" {
						MBRDsk.mbr_partition_4.part_type[0] = byte('P')
					} else {
						var NewEBR EBR
						InitEBR(&NewEBR)
						if Size > int(unsafe.Sizeof(NewEBR)) {
							MBRDsk.mbr_partition_4.part_type[0] = byte('E')
							NewEBR.part_status[0] = 'D'
							NewEBR.part_fit[0] = ' '
							NewEBR.part_start = MBRDsk.mbr_partition_4.part_start
							binary.LittleEndian.PutUint32(NewEBR.part_size, uint32(0))
							binary.LittleEndian.PutUint32(NewEBR.part_next, uint32(0))
							copy(NewEBR.part_name, []byte(""))
							EscribirEBR(Path, NewEBR, Part3Start+Part3Size)
						} else {
							Registrado = false
							Console += "Error, no hay espacio para crear un EBR en la particion Extendida\n"
						}
					}
				default:
					Registrado = false
					Console += "Error, no hay espacio interno suficiente para crear la nueva partición\n"
				}
				if Registrado {
					EscribirMBR(Path, MBRDsk)
					Console += "Partición creada exitosamente\n"
				}
			} else {
				Console += "Error, el disco ya no puede crear mas principales o extendidas\n"
			}
		} else {
			Console += "Error, ya existe una partición extendida\n"
		}
	} else {
		Cabeza := -1
		var PartitionSize int
		if string(MBRDsk.mbr_partition_1.part_status) == "A" && string(MBRDsk.mbr_partition_1.part_type) == "E" {
			Cabeza = int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_1.part_start))
			PartitionSize = int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_1.part_size))
		} else if string(MBRDsk.mbr_partition_2.part_status) == "A" && string(MBRDsk.mbr_partition_2.part_type) == "E" {
			Cabeza = int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_2.part_start))
			PartitionSize = int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_2.part_size))
		} else if string(MBRDsk.mbr_partition_3.part_status) == "A" && string(MBRDsk.mbr_partition_3.part_type) == "E" {
			Cabeza = int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_3.part_start))
			PartitionSize = int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_3.part_size))
		} else if string(MBRDsk.mbr_partition_4.part_status) == "A" && string(MBRDsk.mbr_partition_4.part_type) == "E" {
			Cabeza = int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_4.part_start))
			PartitionSize = int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_4.part_size))
		}
		if Cabeza != -1 {
			var EBRActual EBR
			InitEBR(&EBRActual)
			LeerEBR(Path, &EBRActual, Cabeza)
			EBRNext := int(binary.LittleEndian.Uint32(EBRActual.part_next))
			if string(EBRActual.part_status) == "D" && EBRNext == 0 {
				if Size <= PartitionSize {
					EBRActual.part_status[0] = 'A'
					if strings.ToLower(Fit) == "bf" {
						EBRActual.part_fit[0] = 'B'
					} else if strings.ToLower(Fit) == "ff" {
						EBRActual.part_fit[0] = 'F'
					} else {
						EBRActual.part_fit[0] = 'W'
					}
					binary.LittleEndian.PutUint32(EBRActual.part_start, uint32(Cabeza))
					binary.LittleEndian.PutUint32(EBRActual.part_size, uint32(Size))
					binary.LittleEndian.PutUint32(EBRActual.part_next, uint32(0))
					copy(EBRActual.part_name[:], []byte(Name))
					EscribirEBR(Path, EBRActual, Cabeza)
					Console += "Partición creada exitosamente\n"
				} else {
					Console += "Error, no hay espacio suficiente para crear la partición logica\n"
				}
			} else {
				Actual := Cabeza
				InsertarEN := -1
				Type := ' '
				if strings.ToLower(Fit) == "bf" {
					Type = 'B'
					Menor := PartitionSize
					if string(EBRActual.part_status) == "D" && EBRNext-Cabeza >= Size {
						Menor = EBRNext - Cabeza
						InsertarEN = Cabeza
						Actual = EBRNext
					}
					for Actual != 0 {
						LeerEBR(Path, &EBRActual, Actual)
						EBRNext := int(binary.LittleEndian.Uint32(EBRActual.part_next))
						EBRStart := int(binary.LittleEndian.Uint32(EBRActual.part_start))
						EBRSize := int(binary.LittleEndian.Uint32(EBRActual.part_size))
						if EBRNext != 0 && EBRNext-(EBRStart+EBRSize) >= Size && Menor > EBRNext-(EBRStart+EBRSize) {
							Menor = EBRNext - (EBRStart + EBRSize)
							InsertarEN = EBRStart
						} else if (Cabeza+PartitionSize)-(EBRStart+EBRSize) >= Size && Menor > (Cabeza+PartitionSize)-(EBRStart+EBRSize) {
							InsertarEN = EBRStart
						}
						Actual = EBRNext
					}
				} else if strings.ToLower(Fit) == "ff" {
					Type = 'F'
					Encontrado := false
					if string(EBRActual.part_status) == "D" && EBRNext-Cabeza >= Size {
						InsertarEN = Cabeza
						Actual = EBRNext
						Encontrado = true
					}
					for Actual != 0 && !Encontrado {
						LeerEBR(Path, &EBRActual, Actual)
						EBRNext := int(binary.LittleEndian.Uint32(EBRActual.part_next))
						EBRStart := int(binary.LittleEndian.Uint32(EBRActual.part_start))
						EBRSize := int(binary.LittleEndian.Uint32(EBRActual.part_size))
						if (EBRNext != 0 && EBRNext-(EBRStart+EBRSize) >= Size) || (EBRNext == 0 && (Cabeza+PartitionSize)-(EBRStart+EBRSize) >= Size) {
							Encontrado = true
							InsertarEN = EBRStart
						}
						Actual = EBRNext
					}
				} else {
					Type = 'W'
					Mayor := 0
					if string(EBRActual.part_status) == "D" && EBRNext-Cabeza >= Size {
						Mayor = EBRNext - Cabeza
						InsertarEN = Cabeza
						Actual = EBRNext
					}
					for Actual != 0 {
						LeerEBR(Path, &EBRActual, Actual)
						EBRNext := int(binary.LittleEndian.Uint32(EBRActual.part_next))
						EBRStart := int(binary.LittleEndian.Uint32(EBRActual.part_start))
						EBRSize := int(binary.LittleEndian.Uint32(EBRActual.part_size))
						if EBRNext != 0 && EBRNext-(EBRStart+EBRSize) >= Size && Mayor < EBRNext-(EBRStart+EBRSize) {
							Mayor = EBRNext - (EBRStart + EBRSize)
							InsertarEN = EBRStart
						} else if (Cabeza+PartitionSize)-(EBRStart+EBRSize) >= Size && Mayor < (Cabeza+PartitionSize)-(EBRStart+EBRSize) {
							InsertarEN = EBRStart
						}
						Actual = EBRNext
					}
				}
				if InsertarEN != -1 {
					LeerEBR(Path, &EBRActual, InsertarEN)
					if InsertarEN == Cabeza && string(EBRActual.part_status) == "D" {
						EBRActual.part_status[0] = 'A'
						EBRActual.part_fit[0] = byte(Type)
						binary.LittleEndian.PutUint32(EBRActual.part_size, uint32(Size))
						copy(EBRActual.part_name[:], []byte(Name))
						EscribirEBR(Path, EBRActual, Cabeza)
						Console += "Partición creada exitosamente\n"
					} else {
						EBRNext := int(binary.LittleEndian.Uint32(EBRActual.part_next))
						EBRStart := int(binary.LittleEndian.Uint32(EBRActual.part_size))
						EBRSize := int(binary.LittleEndian.Uint32(EBRActual.part_start))
						var NewEBR EBR
						InitEBR(&NewEBR)
						NewEBR.part_status[0] = 'A'
						NewEBR.part_fit[0] = byte(Type)
						binary.LittleEndian.PutUint32(NewEBR.part_start, uint32(EBRStart+EBRSize))
						binary.LittleEndian.PutUint32(NewEBR.part_size, uint32(Size))
						binary.LittleEndian.PutUint32(NewEBR.part_next, uint32(EBRNext))
						binary.LittleEndian.PutUint32(EBRActual.part_next, uint32(EBRStart+EBRSize))
						copy(NewEBR.part_name[:], []byte(Name))
						EscribirEBR(Path, EBRActual, InsertarEN)
						EscribirEBR(Path, NewEBR, EBRStart+EBRSize)
						Console += "Partición creada exitosamente\n"
					}
				} else {
					Console += "Error, no hay espacio interno suficiente para crear la partición logica\n"
				}
			}
		} else {
			Console += "Error, no hay partición extendida en el disco\n"
		}
	}
}

func MountParticion(Path string, Name string) {
	Mount := false
	for i := 0; i < 10; i++ {
		if !ActivePart[i].Active {
			Mount = true
			var MBRDsk MBR
			InitMBR(&MBRDsk)
			LeerMBR(Path, &MBRDsk)
			NewMount := BuscarParticion(Path, Name)
			if NewMount.Active {
				NewMount.ID = "00" + NewMount.ID + "a"
				NewMount.Path = Path
				ActivePart[i] = NewMount
				Console += "Partición montada con la ID: " + NewMount.ID + "\n"
			}
			break
		} else {
			if Path == ActivePart[i].Path && Name == string(ActivePart[i].part_name) {
				Console += "Error, la partición ya esta montada"
				Mount = true
				break
			}
		}
	}
	if !Mount {
		Console += "Error, ya no se pueden montar mas particiones\n"
	}
}

func BuscarParticion(Path string, Name string) MountPart {
	Name = Init16String(Name)
	var NewMount MountPart
	InitMountPart(&NewMount)
	var MBRDsk MBR
	InitMBR(&MBRDsk)
	LeerMBR(Path, &MBRDsk)
	Cabeza := -1
	if MBRDsk.mbr_partition_1.part_status[0] == 'A' {
		NameA := string(MBRDsk.mbr_partition_1.part_name)
		if NameA == Name {
			NewMount.Active = true
			copy(NewMount.part_name[:], []byte(Name))
			NewMount.part_size = MBRDsk.mbr_partition_1.clone().part_size
			NewMount.part_start = MBRDsk.mbr_partition_1.clone().part_start
			NewMount.ID = "1"
			return NewMount
		}
		if MBRDsk.mbr_partition_1.part_type[0] == 'E' {
			Cabeza = int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_1.part_start))
		}
	}
	if MBRDsk.mbr_partition_2.part_status[0] == 'A' {
		NameA := string(MBRDsk.mbr_partition_2.part_name)
		if NameA == Name {
			NewMount.Active = true
			copy(NewMount.part_name[:], []byte(Name))
			NewMount.part_size = MBRDsk.mbr_partition_2.clone().part_size
			NewMount.part_start = MBRDsk.mbr_partition_2.clone().part_start
			NewMount.ID = "2"
			return NewMount
		}
		if MBRDsk.mbr_partition_2.part_type[0] == 'E' {
			Cabeza = int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_2.part_start))
		}
	}
	if MBRDsk.mbr_partition_3.part_status[0] == 'A' {
		NameA := string(MBRDsk.mbr_partition_3.part_name)
		if NameA == Name {
			NewMount.Active = true
			copy(NewMount.part_name[:], []byte(Name))
			NewMount.part_size = MBRDsk.mbr_partition_3.clone().part_size
			NewMount.part_start = MBRDsk.mbr_partition_3.clone().part_start
			NewMount.ID = "3"
			return NewMount
		}
		if MBRDsk.mbr_partition_3.part_type[0] == 'E' {
			Cabeza = int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_3.part_start))
		}
	}
	if MBRDsk.mbr_partition_4.part_status[0] == 'A' {
		NameA := string(MBRDsk.mbr_partition_4.part_name)
		if NameA == Name {
			NewMount.Active = true
			copy(NewMount.part_name[:], []byte(Name))
			NewMount.part_size = MBRDsk.mbr_partition_4.clone().part_size
			NewMount.part_start = MBRDsk.mbr_partition_4.clone().part_start
			NewMount.ID = "4"
			return NewMount
		}
		if MBRDsk.mbr_partition_4.part_type[0] == 'E' {
			Cabeza = int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_4.part_start))
		}
	}
	if Cabeza != 1 {
		var EBRActual EBR
		InitEBR(&EBRActual)
		LeerEBR(Path, &EBRActual, Cabeza)
		NameA := string(EBRActual.part_name)
		if string(EBRActual.part_status) == "A" && NameA == Name {
			NewMount.Active = true
			copy(NewMount.part_name[:], []byte(Name))
			PartSize := int(binary.LittleEndian.Uint32(EBRActual.part_size))
			PartStart := int(binary.LittleEndian.Uint32(EBRActual.part_start))
			binary.LittleEndian.PutUint32(NewMount.part_size, uint32(PartSize-int(unsafe.Sizeof(EBRActual))))
			binary.LittleEndian.PutUint32(NewMount.part_start, uint32(PartStart+int(unsafe.Sizeof(EBRActual))))
			NewMount.ID = "5"
			return NewMount
		}
		cont := 6
		Next := int(binary.LittleEndian.Uint32(EBRActual.part_next))
		for Next != 0 {
			var EBRNext EBR
			InitEBR(&EBRNext)
			LeerEBR(Path, &EBRNext, Next)
			NameA := string(EBRNext.part_name)
			if NameA == Name {
				NewMount.Active = true
				copy(NewMount.part_name[:], []byte(Name))
				PartSize := int(binary.LittleEndian.Uint32(EBRActual.part_size))
				PartStart := int(binary.LittleEndian.Uint32(EBRActual.part_start))
				binary.LittleEndian.PutUint32(NewMount.part_size, uint32(PartSize-int(unsafe.Sizeof(EBRActual))))
				binary.LittleEndian.PutUint32(NewMount.part_start, uint32(PartStart+int(unsafe.Sizeof(EBRActual))))
				NewMount.ID = strconv.Itoa(cont)
				return NewMount
			}
			EBRActual = EBRNext
			Next = int(binary.LittleEndian.Uint32(EBRActual.part_next))
			cont++
		}
		Console += "Error, no se encontro la partición\n"
	}
	return NewMount
}

func FormatearParticion(ID string) {
	ActiveParticion := -1
	for i := 0; i < 10; i++ {
		if ActivePart[i].Active && ID == ActivePart[i].ID {
			ActiveParticion = i
			break
		}
	}
	if ActiveParticion != -1 {
		var err error
		var NewSuperBloque SuperBloque
		InitSuperBloque(&NewSuperBloque)
		var Raiz TablaInodo
		InitInodo(&Raiz)
		var CarpetaRaiz BloqueCarpeta
		InitBloqueCarpeta(&CarpetaRaiz)

		PartitionSize := int(binary.LittleEndian.Uint32(ActivePart[ActiveParticion].part_size))
		PartitionStart := int(binary.LittleEndian.Uint32(ActivePart[ActiveParticion].part_start))
		ActivePath := ActivePart[ActiveParticion].Path

		N := int(math.Floor(float64(PartitionSize-int(unsafe.Sizeof(NewSuperBloque))) / float64(4+int(unsafe.Sizeof(Raiz))+3*64)))
		FreeInodes := N
		FreeBlocks := 3 * N
		FirstIno := PartitionStart + int(unsafe.Sizeof(NewSuperBloque)) + 4*N
		FirstBlo := PartitionStart + int(unsafe.Sizeof(NewSuperBloque)) + 4*N + N*int(unsafe.Sizeof(Raiz))

		// SB
		binary.LittleEndian.PutUint32(NewSuperBloque.s_filesystem_type, uint32(2))
		binary.LittleEndian.PutUint32(NewSuperBloque.s_inodes_count, uint32(N))
		binary.LittleEndian.PutUint32(NewSuperBloque.s_blocks_count, uint32(3*N))
		binary.LittleEndian.PutUint32(NewSuperBloque.s_free_blocks_count, uint32(FreeBlocks))
		binary.LittleEndian.PutUint32(NewSuperBloque.s_free_inodes_count, uint32(FreeInodes))
		NewSuperBloque.s_mtime, err = (time.Now()).MarshalBinary()
		if err != nil {
			fmt.Println(err)
		}
		binary.LittleEndian.PutUint32(NewSuperBloque.s_mnt_count, uint32(0))
		binary.LittleEndian.PutUint32(NewSuperBloque.s_magic, uint32(0xEF53))
		binary.LittleEndian.PutUint32(NewSuperBloque.s_inode_size, uint32(int(unsafe.Sizeof(Raiz))))
		binary.LittleEndian.PutUint32(NewSuperBloque.s_block_size, uint32(64))
		binary.LittleEndian.PutUint32(NewSuperBloque.s_first_ino, uint32(FirstIno))
		binary.LittleEndian.PutUint32(NewSuperBloque.s_first_blo, uint32(FirstBlo))
		binary.LittleEndian.PutUint32(NewSuperBloque.s_bm_inode_start, uint32(PartitionStart+int(unsafe.Sizeof(NewSuperBloque))))
		binary.LittleEndian.PutUint32(NewSuperBloque.s_bm_block_start, uint32(PartitionStart+int(unsafe.Sizeof(NewSuperBloque))+N))
		binary.LittleEndian.PutUint32(NewSuperBloque.s_inode_start, uint32(PartitionStart+int(unsafe.Sizeof(NewSuperBloque))+4*N))
		binary.LittleEndian.PutUint32(NewSuperBloque.s_block_start, uint32(PartitionStart+int(unsafe.Sizeof(NewSuperBloque))+4*N+N*int(unsafe.Sizeof(Raiz))))

		// Inodo Raiz
		binary.LittleEndian.PutUint32(Raiz.i_uid, uint32(1))
		binary.LittleEndian.PutUint32(Raiz.i_gid, uint32(1))
		binary.LittleEndian.PutUint32(Raiz.i_size, uint32(0))
		Raiz.i_atime, err = (time.Now()).MarshalBinary()
		if err != nil {
			fmt.Println(err)
		}
		Raiz.i_ctime, err = (time.Now()).MarshalBinary()
		if err != nil {
			fmt.Println(err)
		}
		Raiz.i_mtime, err = (time.Now()).MarshalBinary()
		if err != nil {
			fmt.Println(err)
		}
		Raiz.i_type[0] = 0
		binary.LittleEndian.PutUint32(Raiz.i_perm, uint32(770))
		var CarpetaRaizs [16]int
		CarpetaRaizs[0] = 1
		for i := 1; i < 16; i++ {
			CarpetaRaizs[i] = 4294967295
		}
		Raiz.i_block = Int16ArrayToByteArray(CarpetaRaizs)

		EscribirBM(ActivePath, PartitionStart+int(unsafe.Sizeof(NewSuperBloque)), 1)
		EscribirInodo(ActivePath, Raiz, FirstIno)
		FirstIno += int(unsafe.Sizeof(Raiz))
		binary.LittleEndian.PutUint32(NewSuperBloque.s_first_ino, uint32(FirstIno))
		FreeInodes--
		binary.LittleEndian.PutUint32(NewSuperBloque.s_free_inodes_count, uint32(FreeInodes))

		// Bloque Carpeta Raiz
		copy(CarpetaRaiz.b_content[0].b_name[:], []byte(".."))
		binary.LittleEndian.PutUint32(CarpetaRaiz.b_content[0].b_inodo, uint32(0))
		copy(CarpetaRaiz.b_content[1].b_name[:], []byte("."))
		binary.LittleEndian.PutUint32(CarpetaRaiz.b_content[1].b_inodo, uint32(0))
		copy(CarpetaRaiz.b_content[2].b_name[:], []byte("users.txt"))
		binary.LittleEndian.PutUint32(CarpetaRaiz.b_content[2].b_inodo, uint32(1))
		copy(CarpetaRaiz.b_content[3].b_name[:], []byte(""))
		binary.LittleEndian.PutUint32(CarpetaRaiz.b_content[3].b_inodo, uint32(4294967295))

		EscribirBM(ActivePath, PartitionStart+int(unsafe.Sizeof(NewSuperBloque))+N, 1)
		EscribirBloqueCarpeta(ActivePath, CarpetaRaiz, FirstBlo)
		FirstBlo += 64
		binary.LittleEndian.PutUint32(NewSuperBloque.s_first_blo, uint32(FirstBlo))
		FreeBlocks--
		binary.LittleEndian.PutUint32(NewSuperBloque.s_free_blocks_count, uint32(FreeBlocks))

		//Inodo User.txt
		var InodoUser TablaInodo
		InitInodo(&InodoUser)
		binary.LittleEndian.PutUint32(InodoUser.i_uid, uint32(1))
		binary.LittleEndian.PutUint32(InodoUser.i_gid, uint32(1))
		Usuarios := "1, G, root      \n1, U, root      , root      , 123       \n"
		binary.LittleEndian.PutUint32(InodoUser.i_size, uint32(len(Usuarios)))
		InodoUser.i_atime, err = (time.Now()).MarshalBinary()
		if err != nil {
			fmt.Println(err)
		}
		InodoUser.i_ctime, err = (time.Now()).MarshalBinary()
		if err != nil {
			fmt.Println(err)
		}
		InodoUser.i_mtime, err = (time.Now()).MarshalBinary()
		if err != nil {
			fmt.Println(err)
		}
		InodoUser.i_type[0] = 1
		binary.LittleEndian.PutUint32(InodoUser.i_perm, uint32(770))
		CarpetaRaizs[0] = 2
		for i := 1; i < 16; i++ {
			CarpetaRaizs[i] = 4294967295
		}
		InodoUser.i_block = Int16ArrayToByteArray(CarpetaRaizs)

		EscribirBM(ActivePath, PartitionStart+int(unsafe.Sizeof(NewSuperBloque))+1, 1)
		EscribirInodo(ActivePath, InodoUser, FirstIno)
		FirstIno += int(unsafe.Sizeof(Raiz))
		binary.LittleEndian.PutUint32(NewSuperBloque.s_first_ino, uint32(FirstIno))
		FreeInodes--
		binary.LittleEndian.PutUint32(NewSuperBloque.s_free_inodes_count, uint32(FreeInodes))

		//Bloque Archivo Users.txt
		var User_txt BloqueArchivos
		InitBloqueArchivos(&User_txt)
		copy(User_txt.b_content, []byte(Usuarios))
		EscribirBM(ActivePath, PartitionStart+int(unsafe.Sizeof(NewSuperBloque))+N+1, 1)
		EscribirBloqueArchivos(ActivePath, User_txt, FirstBlo)
		FirstBlo += 64
		binary.LittleEndian.PutUint32(NewSuperBloque.s_first_blo, uint32(FirstBlo))
		FreeBlocks--
		binary.LittleEndian.PutUint32(NewSuperBloque.s_free_blocks_count, uint32(FreeBlocks))

		//Escribir SuperBloque
		EscribirSuperBloque(ActivePath, NewSuperBloque, PartitionStart)
	}
}

func Login(User string, Pass string, ID string) {
	ActiveParticion := -1
	for i := 0; i < 10; i++ {
		if ActivePart[i].Active && ID == ActivePart[i].ID {
			ActiveParticion = i
			break
		}
	}
	Activa := ActivePart[ActiveParticion]
	for len(User) < 10 {
		User += " "
	}
	for len(Pass) < 10 {
		Pass += " "
	}
	if ActiveParticion != -1 {

		PartitionStart := int(binary.LittleEndian.Uint32(Activa.part_start))
		var SB SuperBloque
		InitSuperBloque(&SB)
		LeerSuperBloque(Activa.Path, &SB, PartitionStart)
		Temp := Sesion.Active
		Sesion.Active = Activa
		InodeStart := int(binary.LittleEndian.Uint32(SB.s_inode_start))
		var InodeAux TablaInodo
		InitInodo(&InodeAux)
		Usuarios_txt := LeerArchivo(InodeStart + int(unsafe.Sizeof(InodeAux)))
		Sesion.IDU = -1
		Contra := false
		for strings.Index(Usuarios_txt, "\n") != -1 {
			Linea := Usuarios_txt[:strings.Index(Usuarios_txt, "\n")]
			Usuarios_txt = Usuarios_txt[strings.Index(Usuarios_txt, "\n")+1:]
			IDTemp, err := strconv.Atoi(Linea[:strings.Index(Linea, ",")])
			if err != nil {
				fmt.Println(err)
			}
			if IDTemp != 0 {
				Linea = Linea[strings.Index(Linea, ",")+2:]
				Type := Linea[:strings.Index(Linea, ",")]
				Linea = Linea[strings.Index(Linea, ",")+2:]
				if Type == "U" {
					Grupo := Linea[:strings.Index(Linea, ",")]
					Linea = Linea[strings.Index(Linea, ",")+2:]
					TempUser := Linea[:strings.Index(Linea, ",")]
					if TempUser == User {
						Sesion.User = User
						Sesion.IDU = IDTemp
						Sesion.Grupo = Grupo
					}
					Linea = Linea[strings.Index(Linea, ",")+2:]
					if Linea == Pass && Sesion.IDU != -1 {
						Contra = true
					}
				}
			}
		}
		if Sesion.IDU != -1 && Contra {
			Usuarios_txt = LeerArchivo(InodeStart + int(unsafe.Sizeof(InodeAux)))
			for strings.Index(Usuarios_txt, "\n") != -1 {
				Linea := Usuarios_txt[:strings.Index(Usuarios_txt, "\n")]
				Usuarios_txt = Usuarios_txt[strings.Index(Usuarios_txt, "\n")+1:]
				IDTemp, err := strconv.Atoi(Linea[:strings.Index(Linea, ",")])
				if err != nil {
					fmt.Println(err)
				}
				if IDTemp != 0 {
					Linea = Linea[strings.Index(Linea, ",")+2:]
					Type := Linea[:strings.Index(Linea, ",")]
					Linea = Linea[strings.Index(Linea, ",")+2:]
					if Type == "G" {
						if Linea == Sesion.Grupo {
							Sesion.IDG = IDTemp
							Console += "Sesion Iniciada con Exito\n"
							return
						}
					}
				}
			}
			Console += "Error inesperado, grupo no encontrado\n"
		} else {
			Console += "Error, Usuario o Contraseña incorrecta\n"
		}
		Sesion.Active = Temp
		Sesion.User = ""
		Sesion.Grupo = ""
		Sesion.IDU = -1
		Sesion.IDG = -1
	} else {
		Console += "Error, ID de partición no encontrada\n"
	}
}

func AgregarGrupo(Name string) {
	for len(Name) < 10 {
		Name += " "
	}
	PartitionStart := int(binary.LittleEndian.Uint32(Sesion.Active.part_start))
	var SB SuperBloque
	InitSuperBloque(&SB)
	LeerSuperBloque(Sesion.Active.Path, &SB, PartitionStart)
	InodeStart := int(binary.LittleEndian.Uint32(SB.s_inode_start))
	var InodeAux TablaInodo
	InitInodo(&InodeAux)
	Usuarios_txt := LeerArchivo(InodeStart + int(unsafe.Sizeof(InodeAux)))
	IDMayor := 0
	for strings.Index(Usuarios_txt, "\n") != -1 {
		Linea := Usuarios_txt[:strings.Index(Usuarios_txt, "\n")]
		Usuarios_txt = Usuarios_txt[strings.Index(Usuarios_txt, "\n")+1:]
		IDTemp, err := strconv.Atoi(Linea[:strings.Index(Linea, ",")])
		if err != nil {
			fmt.Println(err)
		}
		if IDTemp != 0 {
			Linea = Linea[strings.Index(Linea, ",")+2:]
			Type := Linea[:strings.Index(Linea, ",")]
			Linea = Linea[strings.Index(Linea, ",")+2:]
			if Type == "G" {
				if Linea == Name {
					Console += "Error, ya existe un grupo con este nombre\n"
					return
				} else if IDMayor < IDTemp {
					IDMayor = IDTemp
				}
			}
		}
	}
	ModificarArchivo(InodeStart+int(unsafe.Sizeof(InodeAux)), LeerArchivo(InodeStart+int(unsafe.Sizeof(InodeAux)))+strconv.Itoa(IDMayor+1)+", G, "+Name+"\n")
	Console += "Grupo creado con exito\n"
}

func EliminarGrupo(Name string) {
	for len(Name) < 10 {
		Name += " "
	}
	PartitionStart := int(binary.LittleEndian.Uint32(Sesion.Active.part_start))
	var SB SuperBloque
	InitSuperBloque(&SB)
	LeerSuperBloque(Sesion.Active.Path, &SB, PartitionStart)
	InodeStart := int(binary.LittleEndian.Uint32(SB.s_inode_start))
	var InodeAux TablaInodo
	InitInodo(&InodeAux)
	Usuarios_txt := LeerArchivo(InodeStart + int(unsafe.Sizeof(InodeAux)))
	Aux := ""
	Eliminado := false
	for strings.Index(Usuarios_txt, "\n") != -1 {
		Linea := Usuarios_txt[:strings.Index(Usuarios_txt, "\n")]
		Temp := Linea
		Usuarios_txt = Usuarios_txt[strings.Index(Usuarios_txt, "\n")+1:]
		IDTemp, err := strconv.Atoi(Linea[:strings.Index(Linea, ",")])
		if err != nil {
			fmt.Println(err)
		}
		if IDTemp != 0 {
			Linea = Linea[strings.Index(Linea, ",")+2:]
			Type := Linea[:strings.Index(Linea, ",")]
			Linea = Linea[strings.Index(Linea, ",")+2:]
			if Type == "G" && Linea == Name {
				Eliminado = true
				Aux += "0, G, " + Name + "\n"
			} else {
				Aux += Temp + "\n"
			}
		} else {
			Aux += Temp + "\n"
		}
	}
	if Eliminado {
		ModificarArchivo(InodeStart+int(unsafe.Sizeof(InodeAux)), Aux)
		Console += "Grupo eliminado exitosamente\n"
	} else {
		Console += "Error, grupo no encontrado\n"
	}
}

func AgregarUsuario(User string, Pass string, GRP string) {
	for len(User) < 10 {
		User += " "
	}
	for len(Pass) < 10 {
		Pass += " "
	}
	for len(GRP) < 10 {
		GRP += " "
	}
	PartitionStart := int(binary.LittleEndian.Uint32(Sesion.Active.part_start))
	var SB SuperBloque
	InitSuperBloque(&SB)
	LeerSuperBloque(Sesion.Active.Path, &SB, PartitionStart)
	InodeStart := int(binary.LittleEndian.Uint32(SB.s_inode_start))
	var InodeAux TablaInodo
	InitInodo(&InodeAux)
	Usuarios_txt := LeerArchivo(InodeStart + int(unsafe.Sizeof(InodeAux)))
	IDMayor := 0
	GrupoEncontrado := false
	for strings.Index(Usuarios_txt, "\n") != -1 {
		Linea := Usuarios_txt[:strings.Index(Usuarios_txt, "\n")]
		Usuarios_txt = Usuarios_txt[strings.Index(Usuarios_txt, "\n")+1:]
		IDTemp, err := strconv.Atoi(Linea[:strings.Index(Linea, ",")])
		if err != nil {
			fmt.Println(err)
		}
		if IDTemp != 0 {
			Linea = Linea[strings.Index(Linea, ",")+2:]
			Type := Linea[:strings.Index(Linea, ",")]
			Linea = Linea[strings.Index(Linea, ",")+2:]
			if Type == "G" {
				if Linea == GRP {
					GrupoEncontrado = true
				}
			} else {
				UserName := Linea[:strings.Index(Linea, ",")]
				Linea = Linea[strings.Index(Linea, ",")+2:]
				if User == UserName {
					Console += "Error, ya existe un usuario con este nombre\n"
					return
				} else if IDMayor < IDTemp {
					IDMayor = IDTemp
				}
			}
		}
	}
	if GrupoEncontrado {
		Usuarios_txt = LeerArchivo(InodeStart+int(unsafe.Sizeof(InodeAux))) + strconv.Itoa(IDMayor+1) + ", U, " + User + ", " + GRP + ", " + Pass + "\n"
		ModificarArchivo(InodeStart+int(unsafe.Sizeof(InodeAux)), Usuarios_txt)
		Console += "Usuario creado con exito\n"
	} else {
		Console += "Error, el grupo ingresado no existe\n"
	}
}

func EliminarUsuario(User string) {
	for len(User) < 10 {
		User += " "
	}
	PartitionStart := int(binary.LittleEndian.Uint32(Sesion.Active.part_start))
	var SB SuperBloque
	InitSuperBloque(&SB)
	LeerSuperBloque(Sesion.Active.Path, &SB, PartitionStart)
	InodeStart := int(binary.LittleEndian.Uint32(SB.s_inode_start))
	var InodeAux TablaInodo
	InitInodo(&InodeAux)
	Usuarios_txt := LeerArchivo(InodeStart + int(unsafe.Sizeof(InodeAux)))
	Aux := ""
	Eliminado := false
	for strings.Index(Usuarios_txt, "\n") != -1 {
		Linea := Usuarios_txt[:strings.Index(Usuarios_txt, "\n")]
		Temp := Linea
		Usuarios_txt = Usuarios_txt[strings.Index(Usuarios_txt, "\n")+1:]
		IDTemp, err := strconv.Atoi(Linea[:strings.Index(Linea, ",")])
		if err != nil {
			fmt.Println(err)
		}
		if IDTemp != 0 {
			Linea = Linea[strings.Index(Linea, ",")+2:]
			Type := Linea[:strings.Index(Linea, ",")]
			Linea = Linea[strings.Index(Linea, ",")+2:]
			if Type == "U" {
				Name := Linea[:strings.Index(Linea, ",")]
				Linea = Linea[strings.Index(Linea, ",")+2:]
				if User == Name {
					Eliminado = true
					Aux += "0, U, " + User + ", " + Linea + "\n"
				} else {
					Aux += Temp + "\n"
				}
			} else {
				Aux += Temp + "\n"
			}
		} else {
			Aux += Temp + "\n"
		}
	}
	if Eliminado {
		ModificarArchivo(InodeStart+int(unsafe.Sizeof(InodeAux)), Aux)
		Console += "Usuario eliminado exitosamente\n"
	} else {
		Console += "Error, grupo no encontrado\n"
	}
}

func MKFILE(Path string, R bool, Size int, Cont string) {
	PartitionStart := int(binary.LittleEndian.Uint32(Sesion.Active.part_start))
	var SB SuperBloque
	InitSuperBloque(&SB)
	LeerSuperBloque(Sesion.Active.Path, &SB, PartitionStart)
	InodeStart := int(binary.LittleEndian.Uint32(SB.s_inode_start))
	if Path[0] == '/' {
		InodeActual := int(binary.LittleEndian.Uint32(SB.s_inode_start))
		Path = Path[1:]
		SubPath := "/"
		Error := false
		for strings.Index(Path, "/") != -1 {
			Name := Path[:strings.Index(Path, "/")]
			Path = Path[strings.Index(Path, "/")+1:]
			NextInodo := BuscarCarpeta(InodeActual, Name)
			var InodoAux TablaInodo
			InitInodo(&InodoAux)
			if NextInodo != 4294967295 {
				InodeActual = InodeStart + NextInodo*int(unsafe.Sizeof(InodoAux))
			} else {
				if R {
					if ComprobarPermisos(InodeActual, false, true, false) {
						CrearCarpeta(InodeActual, Name)
						NextInodo = BuscarCarpeta(InodeActual, Name)
						InodeActual = InodeStart + NextInodo*int(unsafe.Sizeof(InodoAux))
					} else {
						Console += "Error, no tiene permisos para crear una carpeta en \"" + SubPath + "\"\n"
						Error = true
					}
				} else {
					Console += "Error, el path debe iniciar en la carpeta raiz\n"
					Error = true
				}
			}
			SubPath = Name + "/"
		}
		if !Error {
			if BuscarCarpeta(InodeActual, Path) == 4294967295 {
				if ComprobarPermisos(InodeActual, false, true, false) {
					CrearArchivo(InodeActual, Path, Size, Cont)
					Console += "Archivo creado con exito\n"
				} else {
					Console += "Error, no tiene permisos para crear una archivo en \"" + SubPath + "\"\n"
					Error = true
				}
			} else {
				Console += "Error, ya hay un archivo con ese nombre en la carpeta indicada\n"
			}
		}
	} else {
		Console += "Error, el path debe iniciar en la carpeta raiz\n"
	}
}

//Funciones de manejo de carpetas y archivos

func LeerArchivo(InodoInit int) string {
	Contenido := ""
	Activa := Sesion.Active
	PartitionStart := int(binary.LittleEndian.Uint32(Activa.part_start))
	var SB SuperBloque
	InitSuperBloque(&SB)
	LeerSuperBloque(Activa.Path, &SB, PartitionStart)
	BlockStart := int(binary.LittleEndian.Uint32(SB.s_block_start))
	var Aux TablaInodo
	InitInodo(&Aux)
	LeerInodo(Activa.Path, &Aux, InodoInit)
	TempBlocks := ByteArrayToInt16Array(Aux.i_block)
	for i := 0; i < 16; i++ {
		if TempBlocks[i] != 4294967295 {
			var Temp BloqueArchivos
			InitBloqueArchivos(&Temp)
			LeerBloqueArchivos(Activa.Path, &Temp, BlockStart+(TempBlocks[i]-1)*64)
			TempContent := string(Temp.b_content)
			for j := 0; j < len(TempContent); j++ {
				if TempContent[j] != 0 {
					Contenido += string(TempContent[j])
				}
			}
		}
	}
	return Contenido
}

func ModificarArchivo(InodoInit int, Texto string) {
	PartitionStart := int(binary.LittleEndian.Uint32(Sesion.Active.part_start))
	var SB SuperBloque
	InitSuperBloque(&SB)
	LeerSuperBloque(Sesion.Active.Path, &SB, PartitionStart)
	//InodeStart := int(binary.LittleEndian.Uint32(SB.s_inode_start))
	BMBlockStart := int(binary.LittleEndian.Uint32(SB.s_bm_block_start))
	var Archivo TablaInodo
	InitInodo(&Archivo)
	LeerInodo(Sesion.Active.Path, &Archivo, InodoInit)
	binary.LittleEndian.PutUint32(Archivo.i_size, uint32(len(Texto)))
	EscribirInodo(Sesion.Active.Path, Archivo, InodoInit)
	TempBlocks := ByteArrayToInt16Array(Archivo.i_block)
	for i := 0; i < 16; i++ {
		if TempBlocks[i] != 4294967295 {
			if Texto != "" {
				Texto = ModificarBloqueArchivo(TempBlocks[i], Texto)
			} else {
				EscribirBM(Sesion.Active.Path, BMBlockStart+TempBlocks[i]-1, 0)
				TempBlocks[i] = 4294967295
				Archivo.i_block = Int16ArrayToByteArray(TempBlocks)
				EscribirInodo(Sesion.Active.Path, Archivo, InodoInit)
				binary.LittleEndian.PutUint32(SB.s_free_blocks_count, uint32(int(binary.LittleEndian.Uint32(SB.s_free_blocks_count))+1))
				binary.LittleEndian.PutUint32(SB.s_first_blo, uint32(int(binary.LittleEndian.Uint32(SB.s_free_blocks_count)))+uint32(FirstFreeBlock()*64))
				EscribirSuperBloque(Sesion.Active.Path, SB, PartitionStart)
			}
		} else if Texto != "" {
			Nuevo := FirstFreeBlock()
			EscribirBM(Sesion.Active.Path, BMBlockStart+Nuevo, 1)
			var NuevoBloque BloqueArchivos
			InitBloqueArchivos(&NuevoBloque)
			copy(NuevoBloque.b_content, []byte(""))
			EscribirBloqueArchivos(Sesion.Active.Path, NuevoBloque, Nuevo*64)
			LeerInodo(Sesion.Active.Path, &Archivo, InodoInit)
			TempBlocks[i] = Nuevo + 1
			Archivo.i_block = Int16ArrayToByteArray(TempBlocks)
			EscribirInodo(Sesion.Active.Path, Archivo, InodoInit)
			binary.LittleEndian.PutUint32(SB.s_free_blocks_count, uint32(int(binary.LittleEndian.Uint32(SB.s_free_blocks_count))-1))
			binary.LittleEndian.PutUint32(SB.s_first_blo, uint32(int(binary.LittleEndian.Uint32(SB.s_free_blocks_count)))+uint32(FirstFreeBlock()*64))
			EscribirSuperBloque(Sesion.Active.Path, SB, PartitionStart)
			Texto = ModificarBloqueArchivo(TempBlocks[i], Texto)
		}
	}
}

func ModificarBloqueArchivo(NumBloque int, Texto string) string {
	PartitionStart := int(binary.LittleEndian.Uint32(Sesion.Active.part_start))
	var SB SuperBloque
	InitSuperBloque(&SB)
	LeerSuperBloque(Sesion.Active.Path, &SB, PartitionStart)
	BlockStart := int(binary.LittleEndian.Uint32(SB.s_block_start))
	var Archivo BloqueArchivos
	InitBloqueArchivos(&Archivo)
	LeerBloqueArchivos(Sesion.Active.Path, &Archivo, BlockStart+(NumBloque-1)*64)
	if len(Texto) > 64 {
		copy(Archivo.b_content, []byte(Texto[0:64]))
		Texto = Texto[64:]
	} else {
		copy(Archivo.b_content, []byte(Texto))
		Texto = ""
	}
	EscribirBloqueArchivos(Sesion.Active.Path, Archivo, BlockStart+(NumBloque-1)*64)
	return Texto
}

func BuscarCarpeta(InodoActual int, Name string) int {
	PartitionStart := int(binary.LittleEndian.Uint32(Sesion.Active.part_start))
	var SB SuperBloque
	InitSuperBloque(&SB)
	LeerSuperBloque(Sesion.Active.Path, &SB, PartitionStart)
	BlockStart := int(binary.LittleEndian.Uint32(SB.s_block_start))
	var InodoPadre TablaInodo
	InitInodo(&InodoPadre)
	LeerInodo(Sesion.Active.Path, &InodoPadre, InodoActual)
	Blocks := ByteArrayToInt16Array(InodoPadre.i_block)
	for i := 0; i < 16; i++ {
		if Blocks[i] != 4294967295 {
			Temp := BuscarCarpetaEnBloque(BlockStart+(Blocks[i]-1)*64, Name)
			if Temp != 4294967295 {
				return Temp
			}
		}
	}
	return 4294967295
}

func BuscarCarpetaEnBloque(BlockActual int, Name string) int {
	PartitionStart := int(binary.LittleEndian.Uint32(Sesion.Active.part_start))
	var SB SuperBloque
	InitSuperBloque(&SB)
	LeerSuperBloque(Sesion.Active.Path, &SB, PartitionStart)
	var Actual BloqueCarpeta
	InitBloqueCarpeta(&Actual)
	LeerBloqueCarpeta(Sesion.Active.Path, &Actual, BlockActual)
	for i := 0; i < 4; i++ {
		TempInode := int(binary.LittleEndian.Uint32(Actual.b_content[i].b_inodo))
		if TempInode != 4294967295 {
			if Init12String(Name) == string(Actual.b_content[i].b_name) {
				return TempInode
			}
		}
	}
	return 4294967295
}

func CrearCarpeta(InodoActual int, Name string) bool {
	PartitionStart := int(binary.LittleEndian.Uint32(Sesion.Active.part_start))
	var SB SuperBloque
	InitSuperBloque(&SB)
	LeerSuperBloque(Sesion.Active.Path, &SB, PartitionStart)
	BMBlockStart := int(binary.LittleEndian.Uint32(SB.s_bm_block_start))
	BlockStart := int(binary.LittleEndian.Uint32(SB.s_block_start))
	var InodoPadre TablaInodo
	InitInodo(&InodoPadre)
	LeerInodo(Sesion.Active.Path, &InodoPadre, InodoActual)
	Blocks := ByteArrayToInt16Array(InodoPadre.i_block)
	for i := 0; i < 16; i++ {
		if Blocks[i] == 4294967295 {
			Nuevo := FirstFreeBlock()
			EscribirBM(Sesion.Active.Path, BMBlockStart+Nuevo, 1)
			var NuevoBloque BloqueCarpeta
			InitBloqueCarpeta(&NuevoBloque)
			for j := 0; j < 4; j++ {
				copy(NuevoBloque.b_content[j].b_name[:], []byte(""))
				binary.LittleEndian.PutUint32(NuevoBloque.b_content[j].b_inodo, uint32(4294967295))
			}
			EscribirBloqueCarpeta(Sesion.Active.Path, NuevoBloque, BlockStart+Nuevo*64)
			Blocks[i] = Nuevo + 1
			InodoPadre.i_block = Int16ArrayToByteArray(Blocks)
			EscribirInodo(Sesion.Active.Path, InodoPadre, InodoActual)
			binary.LittleEndian.PutUint32(SB.s_first_blo, uint32(BlockStart+FirstFreeBlock()*64))
			EscribirSuperBloque(Sesion.Active.Path, SB, PartitionStart)
		}
		if CrearCarpetaEnBloque(BlockStart+(Blocks[i]-1)*64, Name, InodoActual) {
			return true
		}
	}
	return false
}

func CrearCarpetaEnBloque(Bloque int, Name string, Anterior int) bool {
	PartitionStart := int(binary.LittleEndian.Uint32(Sesion.Active.part_start))
	var SB SuperBloque
	InitSuperBloque(&SB)
	LeerSuperBloque(Sesion.Active.Path, &SB, PartitionStart)
	BMInodeStart := int(binary.LittleEndian.Uint32(SB.s_bm_inode_start))
	InodeStart := int(binary.LittleEndian.Uint32(SB.s_inode_start))
	BMBlockStart := int(binary.LittleEndian.Uint32(SB.s_bm_block_start))
	BlockStart := int(binary.LittleEndian.Uint32(SB.s_block_start))
	var Actual BloqueCarpeta
	InitBloqueCarpeta(&Actual)
	LeerBloqueCarpeta(Sesion.Active.Path, &Actual, Bloque)
	for i := 0; i < 4; i++ {
		if int(binary.LittleEndian.Uint32(Actual.b_content[i].b_inodo)) == 4294967295 {
			Nuevo := FirstFreeInode()
			EscribirBM(Sesion.Active.Path, BMInodeStart+Nuevo, 1)
			var NuevoInodo TablaInodo
			InitInodo(&NuevoInodo)
			binary.LittleEndian.PutUint32(NuevoInodo.i_uid, uint32(Sesion.IDU))
			binary.LittleEndian.PutUint32(NuevoInodo.i_gid, uint32(Sesion.IDG))
			binary.LittleEndian.PutUint32(NuevoInodo.i_size, uint32(0))
			var Blocks [16]int
			Blocks[0] = FirstFreeBlock() + 1
			for j := 1; j < 16; j++ {
				Blocks[j] = 4294967295
			}
			NuevoInodo.i_block = Int16ArrayToByteArray(Blocks)
			NuevoInodo.i_type[0] = 0
			binary.LittleEndian.PutUint32(NuevoInodo.i_perm, uint32(664))
			EscribirInodo(Sesion.Active.Path, NuevoInodo, InodeStart+Nuevo*int(unsafe.Sizeof(NuevoInodo)))

			binary.LittleEndian.PutUint32(Actual.b_content[i].b_inodo, uint32(Nuevo))
			copy(Actual.b_content[i].b_name[:], []byte(Name))

			NuevoB := FirstFreeBlock()
			EscribirBM(Sesion.Active.Path, BMBlockStart+NuevoB, 1)
			var NuevoBloque BloqueCarpeta
			InitBloqueCarpeta(&NuevoBloque)
			binary.LittleEndian.PutUint32(NuevoBloque.b_content[0].b_inodo, uint32(Nuevo))
			copy(NuevoBloque.b_content[0].b_name[:], []byte("."))
			binary.LittleEndian.PutUint32(NuevoBloque.b_content[1].b_inodo, uint32((Anterior-InodeStart)/int(unsafe.Sizeof(NuevoInodo))))
			copy(NuevoBloque.b_content[1].b_name[:], []byte(".."))
			binary.LittleEndian.PutUint32(NuevoBloque.b_content[2].b_inodo, uint32(4294967295))
			copy(NuevoBloque.b_content[2].b_name[:], []byte(""))
			binary.LittleEndian.PutUint32(NuevoBloque.b_content[3].b_inodo, uint32(4294967295))
			copy(NuevoBloque.b_content[3].b_name[:], []byte(""))
			EscribirBloqueCarpeta(Sesion.Active.Path, NuevoBloque, BlockStart+NuevoB*64)

			EscribirBloqueCarpeta(Sesion.Active.Path, Actual, Bloque)

			binary.LittleEndian.PutUint32(SB.s_first_ino, uint32(InodeStart+FirstFreeInode()*int(unsafe.Sizeof(NuevoInodo))))
			binary.LittleEndian.PutUint32(SB.s_first_blo, uint32(BlockStart+FirstFreeBlock()*64))
			EscribirSuperBloque(Sesion.Active.Path, SB, PartitionStart)
			return true
		}
	}
	return false
}

func CrearArchivo(InodoActual int, Name string, Size int, Cont string) bool {
	PartitionStart := int(binary.LittleEndian.Uint32(Sesion.Active.part_start))
	var SB SuperBloque
	InitSuperBloque(&SB)
	LeerSuperBloque(Sesion.Active.Path, &SB, PartitionStart)
	BMBlockStart := int(binary.LittleEndian.Uint32(SB.s_bm_block_start))
	BlockStart := int(binary.LittleEndian.Uint32(SB.s_block_start))
	var InodoPadre TablaInodo
	InitInodo(&InodoPadre)
	LeerInodo(Sesion.Active.Path, &InodoPadre, InodoActual)
	Blocks := ByteArrayToInt16Array(InodoPadre.i_block)
	for i := 0; i < 16; i++ {
		if Blocks[i] == 4294967295 {
			Nuevo := FirstFreeBlock()
			EscribirBM(Sesion.Active.Path, BMBlockStart+Nuevo, 1)
			var NuevoBloque BloqueCarpeta
			InitBloqueCarpeta(&NuevoBloque)
			for j := 0; j < 4; j++ {
				copy(NuevoBloque.b_content[j].b_name[:], []byte(""))
				binary.LittleEndian.PutUint32(NuevoBloque.b_content[j].b_inodo, uint32(4294967295))
			}
			EscribirBloqueCarpeta(Sesion.Active.Path, NuevoBloque, BlockStart+Nuevo*64)
			Blocks[i] = Nuevo + 1
			InodoPadre.i_block = Int16ArrayToByteArray(Blocks)
			EscribirInodo(Sesion.Active.Path, InodoPadre, InodoActual)
			binary.LittleEndian.PutUint32(SB.s_first_blo, uint32(BlockStart+FirstFreeBlock()*64))
			EscribirSuperBloque(Sesion.Active.Path, SB, PartitionStart)
		}
		if CrearArchivoEnBloque(BlockStart+(Blocks[i]-1)*64, Name, InodoActual, Size, Cont) {
			return true
		}
	}
	return false
}

func CrearArchivoEnBloque(Bloque int, Name string, Anterior int, Size int, Cont string) bool {
	PartitionStart := int(binary.LittleEndian.Uint32(Sesion.Active.part_start))
	var SB SuperBloque
	InitSuperBloque(&SB)
	LeerSuperBloque(Sesion.Active.Path, &SB, PartitionStart)
	BMInodeStart := int(binary.LittleEndian.Uint32(SB.s_bm_inode_start))
	InodeStart := int(binary.LittleEndian.Uint32(SB.s_inode_start))
	var Actual BloqueCarpeta
	InitBloqueCarpeta(&Actual)
	LeerBloqueCarpeta(Sesion.Active.Path, &Actual, Bloque)
	for i := 0; i < 4; i++ {
		if int(binary.LittleEndian.Uint32(Actual.b_content[i].b_inodo)) == 4294967295 {
			Nuevo := FirstFreeInode()
			EscribirBM(Sesion.Active.Path, BMInodeStart+Nuevo, 1)
			var NuevoInodo TablaInodo
			InitInodo(&NuevoInodo)
			binary.LittleEndian.PutUint32(NuevoInodo.i_uid, uint32(Sesion.IDU))
			binary.LittleEndian.PutUint32(NuevoInodo.i_gid, uint32(Sesion.IDG))
			binary.LittleEndian.PutUint32(NuevoInodo.i_size, uint32(0))
			var Blocks [16]int
			for j := 0; j < 16; j++ {
				Blocks[j] = 4294967295
			}
			NuevoInodo.i_block = Int16ArrayToByteArray(Blocks)
			NuevoInodo.i_type[0] = 1
			binary.LittleEndian.PutUint32(NuevoInodo.i_perm, uint32(664))
			EscribirInodo(Sesion.Active.Path, NuevoInodo, InodeStart+Nuevo*int(unsafe.Sizeof(NuevoInodo)))

			binary.LittleEndian.PutUint32(Actual.b_content[i].b_inodo, uint32(Nuevo))
			copy(Actual.b_content[i].b_name[:], []byte(Name))
			EscribirBloqueCarpeta(Sesion.Active.Path, Actual, Bloque)

			binary.LittleEndian.PutUint32(SB.s_first_ino, uint32(InodeStart+FirstFreeInode()*int(unsafe.Sizeof(NuevoInodo))))
			EscribirSuperBloque(Sesion.Active.Path, SB, PartitionStart)

			if Cont != "" {
				file, err := os.Open(Cont)
				if err != nil {
					fmt.Println(err)
				}
				defer file.Close()
				content := ""
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					content += scanner.Text() + "\n"
				}
				ModificarArchivo(InodeStart+Nuevo*int(unsafe.Sizeof(NuevoInodo)), content[:len(content)-1])
			} else if Size != 0 {
				TempChars := make([]byte, Size)
				copy(TempChars, []byte(""))
				ModificarArchivo(InodeStart+Nuevo*int(unsafe.Sizeof(NuevoInodo)), string(TempChars))
			}
			return true
		}
	}
	return false
}

func ComprobarPermisos(InodoActual int, Lectura bool, Escritura bool, Ejecucion bool) bool {
	if Sesion.User == "root      " {
		return true
	}
	var Actual TablaInodo
	InitInodo(&Actual)
	LeerInodo(Sesion.Active.Path, &Actual, InodoActual)
	Perms := int(binary.LittleEndian.Uint32(Actual.i_perm))
	var U, G, O int
	U = Perms
	O = U % 10
	U = (U - O) / 10
	G = U % 10
	U = (U - G) / 10
	var Permiso int
	UID := int(binary.LittleEndian.Uint32(Actual.i_uid))
	GID := int(binary.LittleEndian.Uint32(Actual.i_gid))
	if UID == Sesion.IDU {
		Permiso = U
	} else if GID == Sesion.IDG {
		Permiso = G
	} else {
		Permiso = O
	}
	var ALectura, AEscritura, AEjecucion bool
	switch Permiso {
	case 1:
		ALectura = false
		AEscritura = false
		AEjecucion = true
	case 2:
		ALectura = false
		AEscritura = true
		AEjecucion = false
	case 3:
		ALectura = false
		AEscritura = true
		AEjecucion = true
	case 4:
		ALectura = true
		AEscritura = false
		AEjecucion = false
	case 5:
		ALectura = true
		AEscritura = false
		AEjecucion = true
	case 6:
		Lectura = true
		AEscritura = true
		AEjecucion = false
	case 7:
		Lectura = true
		AEscritura = true
		AEjecucion = true
	default:
		ALectura = false
		AEscritura = false
		AEjecucion = false
	}
	if Lectura && !ALectura || Escritura && !AEscritura || Ejecucion && !AEjecucion {
		return false
	} else {
		return true
	}
}

func FirstFreeBlock() int {
	PartitionStart := int(binary.LittleEndian.Uint32(Sesion.Active.part_start))
	var SB SuperBloque
	InitSuperBloque(&SB)
	LeerSuperBloque(Sesion.Active.Path, &SB, PartitionStart)
	BlocksCount := int(binary.LittleEndian.Uint32(SB.s_blocks_count))
	BMBlockStart := int(binary.LittleEndian.Uint32(SB.s_bm_block_start))
	for i := 0; i < BlocksCount; i++ {
		if LeerBM(Sesion.Active.Path, BMBlockStart+i) == 0 {
			return i
		}
	}
	return -1
}

func FirstFreeInode() int {
	PartitionStart := int(binary.LittleEndian.Uint32(Sesion.Active.part_start))
	var SB SuperBloque
	InitSuperBloque(&SB)
	LeerSuperBloque(Sesion.Active.Path, &SB, PartitionStart)
	InodesCount := int(binary.LittleEndian.Uint32(SB.s_inodes_count))
	BMInodeStart := int(binary.LittleEndian.Uint32(SB.s_bm_inode_start))
	for i := 0; i < InodesCount; i++ {
		if LeerBM(Sesion.Active.Path, BMInodeStart+i) == 0 {
			return i
		}
	}
	return -1
}

//Funciones auxiliares

func Valor(Linea *string) string {
	Val := ""
	if (*Linea)[0] != '"' {
		if strings.Index(*Linea, " ") != -1 && (strings.Index(*Linea, "=") == -1 || strings.Index(*Linea, " ") < strings.Index(*Linea, "=")) {
			Val = (*Linea)[0:strings.Index(*Linea, " ")]
			*Linea = (*Linea)[strings.Index(*Linea, " ")+1 : len((*Linea))]
		} else if strings.Index(*Linea, "=") != -1 {
			Val = (*Linea)[0:strings.Index(*Linea, "=")]
			*Linea = (*Linea)[strings.Index(*Linea, "=")+1 : len((*Linea))]
		}
	} else {
		*Linea = (*Linea)[1:len((*Linea))]
		Val = (*Linea)[0:strings.Index(*Linea, "\"")]
		*Linea = (*Linea)[strings.Index(*Linea, "\"")+1 : len((*Linea))]
	}
	for len(*Linea) > 0 && (*Linea)[0] == ' ' {
		*Linea = (*Linea)[1:len((*Linea))]
	}
	return Val
}

func CrearPath(Path string) {
	SubPath := ""
	for strings.Index(Path, "/") != -1 {
		SubPath += Path[0 : strings.Index(Path, "/")+1]
		Path = Path[strings.Index(Path, "/")+1:]
		os.Mkdir(SubPath, 0777)
	}
}

//Inits, Escribir y Leer Structs

func Inicializacion() {
	for i := 0; i < 10; i++ {
		ActivePart[i].Active = false
	}
	Sesion.User = ""
	Sesion.Grupo = ""
	Sesion.IDU = -1
	Sesion.IDG = -1
	Sesion.Active.Active = false
	Console = ""
	rand.Seed(time.Now().UnixNano())
}

func Init16String(Original string) string {
	Temp := make([]byte, 16)
	copy(Temp, []byte(Original))
	return string(Temp)
}

func Init12String(Original string) string {
	Temp := make([]byte, 12)
	copy(Temp, []byte(Original))
	return string(Temp)
}

func InitMBR(NewMBR *MBR) {
	var err error
	NewMBR.mbr_tamano = make([]byte, 4)
	NewMBR.mbr_fecha_creacion, err = (time.Now()).MarshalBinary()
	if err != nil {
		fmt.Println(err)
	}
	NewMBR.mbr_dsk_signature = make([]byte, 4)
	NewMBR.mbr_dsk_fit = make([]byte, 1)
	NewMBR.mbr_partition_1.part_status = make([]byte, 1)
	NewMBR.mbr_partition_1.part_type = make([]byte, 1)
	NewMBR.mbr_partition_1.part_fit = make([]byte, 1)
	NewMBR.mbr_partition_1.part_start = make([]byte, 4)
	NewMBR.mbr_partition_1.part_size = make([]byte, 4)
	NewMBR.mbr_partition_1.part_name = make([]byte, 16)
	NewMBR.mbr_partition_2.part_status = make([]byte, 1)
	NewMBR.mbr_partition_2.part_type = make([]byte, 1)
	NewMBR.mbr_partition_2.part_fit = make([]byte, 1)
	NewMBR.mbr_partition_2.part_start = make([]byte, 4)
	NewMBR.mbr_partition_2.part_size = make([]byte, 4)
	NewMBR.mbr_partition_2.part_name = make([]byte, 16)
	NewMBR.mbr_partition_3.part_status = make([]byte, 1)
	NewMBR.mbr_partition_3.part_type = make([]byte, 1)
	NewMBR.mbr_partition_3.part_fit = make([]byte, 1)
	NewMBR.mbr_partition_3.part_start = make([]byte, 4)
	NewMBR.mbr_partition_3.part_size = make([]byte, 4)
	NewMBR.mbr_partition_3.part_name = make([]byte, 16)
	NewMBR.mbr_partition_4.part_status = make([]byte, 1)
	NewMBR.mbr_partition_4.part_type = make([]byte, 1)
	NewMBR.mbr_partition_4.part_fit = make([]byte, 1)
	NewMBR.mbr_partition_4.part_start = make([]byte, 4)
	NewMBR.mbr_partition_4.part_size = make([]byte, 4)
	NewMBR.mbr_partition_4.part_name = make([]byte, 16)
}

func EscribirMBR(Path string, MBRActual MBR) {
	archivo, err := os.OpenFile(Path, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println(err)
	}
	defer archivo.Close()
	archivo.Seek(0, 0)
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_tamano)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_fecha_creacion)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_dsk_signature)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_dsk_fit)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_1.part_status)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_1.part_type)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_1.part_fit)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_1.part_start)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_1.part_size)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_1.part_name)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_2.part_status)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_2.part_type)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_2.part_fit)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_2.part_start)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_2.part_size)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_2.part_name)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_3.part_status)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_3.part_type)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_3.part_fit)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_3.part_start)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_3.part_size)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_3.part_name)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_4.part_status)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_4.part_type)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_4.part_fit)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_4.part_start)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_4.part_size)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &MBRActual.mbr_partition_4.part_name)
	if err != nil {
		fmt.Println(err)
	}
}

func LeerMBR(Path string, MBRActual *MBR) {
	archivo, err := os.Open(Path)
	if err != nil {
		fmt.Println(err)
	}
	defer archivo.Close()
	archivo.Seek(0, 0)
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_tamano)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_fecha_creacion)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_dsk_signature)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_dsk_fit)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_1.part_status)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_1.part_type)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_1.part_fit)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_1.part_start)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_1.part_size)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_1.part_name)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_2.part_status)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_2.part_type)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_2.part_fit)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_2.part_start)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_2.part_size)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_2.part_name)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_3.part_status)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_3.part_type)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_3.part_fit)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_3.part_start)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_3.part_size)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_3.part_name)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_4.part_status)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_4.part_type)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_4.part_fit)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_4.part_start)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_4.part_size)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, MBRActual.mbr_partition_4.part_name)
	if err != nil {
		fmt.Println(err)
	}
}

func InitEBR(NewEBR *EBR) {
	NewEBR.part_status = make([]byte, 1)
	NewEBR.part_fit = make([]byte, 1)
	NewEBR.part_start = make([]byte, 4)
	NewEBR.part_size = make([]byte, 4)
	NewEBR.part_next = make([]byte, 4)
	NewEBR.part_name = make([]byte, 16)
}

func EscribirEBR(Path string, EBRActual EBR, index int) {
	archivo, err := os.OpenFile(Path, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println(err)
	}
	defer archivo.Close()
	archivo.Seek(int64(index), 0)
	err = binary.Write(archivo, binary.LittleEndian, &EBRActual.part_status)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &EBRActual.part_fit)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &EBRActual.part_start)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &EBRActual.part_size)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &EBRActual.part_next)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &EBRActual.part_name)
	if err != nil {
		fmt.Println(err)
	}
}

func LeerEBR(Path string, EBRActual *EBR, index int) {
	archivo, err := os.Open(Path)
	if err != nil {
		fmt.Println(err)
	}
	defer archivo.Close()
	archivo.Seek(int64(index), 0)
	err = binary.Read(archivo, binary.LittleEndian, EBRActual.part_status)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, EBRActual.part_fit)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, EBRActual.part_start)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, EBRActual.part_size)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, EBRActual.part_next)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, EBRActual.part_name)
	if err != nil {
		fmt.Println(err)
	}
}

func InitMountPart(NewMountPart *MountPart) {
	NewMountPart.part_start = make([]byte, 4)
	NewMountPart.part_size = make([]byte, 4)
	NewMountPart.part_name = make([]byte, 16)
}

func InitSuperBloque(NewSuperBloque *SuperBloque) {
	var err error
	NewSuperBloque.s_filesystem_type = make([]byte, 4)
	NewSuperBloque.s_inodes_count = make([]byte, 4)
	NewSuperBloque.s_blocks_count = make([]byte, 4)
	NewSuperBloque.s_free_blocks_count = make([]byte, 4)
	NewSuperBloque.s_free_inodes_count = make([]byte, 4)
	NewSuperBloque.s_mtime, err = (time.Now()).MarshalBinary()
	if err != nil {
		fmt.Println(err)
	}
	NewSuperBloque.s_mnt_count = make([]byte, 4)
	NewSuperBloque.s_magic = make([]byte, 4)
	NewSuperBloque.s_inode_size = make([]byte, 4)
	NewSuperBloque.s_block_size = make([]byte, 4)
	NewSuperBloque.s_first_ino = make([]byte, 4)
	NewSuperBloque.s_first_blo = make([]byte, 4)
	NewSuperBloque.s_bm_inode_start = make([]byte, 4)
	NewSuperBloque.s_bm_block_start = make([]byte, 4)
	NewSuperBloque.s_inode_start = make([]byte, 4)
	NewSuperBloque.s_block_start = make([]byte, 4)
}

func EscribirSuperBloque(Path string, SuperBloqueActual SuperBloque, index int) {
	archivo, err := os.OpenFile(Path, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println(err)
	}
	defer archivo.Close()
	archivo.Seek(int64(index), 0)
	err = binary.Write(archivo, binary.LittleEndian, &SuperBloqueActual.s_filesystem_type)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &SuperBloqueActual.s_inodes_count)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &SuperBloqueActual.s_blocks_count)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &SuperBloqueActual.s_free_blocks_count)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &SuperBloqueActual.s_free_inodes_count)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &SuperBloqueActual.s_mtime)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &SuperBloqueActual.s_mnt_count)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &SuperBloqueActual.s_magic)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &SuperBloqueActual.s_inode_size)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &SuperBloqueActual.s_block_size)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &SuperBloqueActual.s_first_ino)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &SuperBloqueActual.s_first_blo)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &SuperBloqueActual.s_bm_inode_start)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &SuperBloqueActual.s_bm_block_start)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &SuperBloqueActual.s_inode_start)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &SuperBloqueActual.s_block_start)
	if err != nil {
		fmt.Println(err)
	}
}

func LeerSuperBloque(Path string, SuperBloqueActual *SuperBloque, index int) {
	archivo, err := os.Open(Path)
	if err != nil {
		fmt.Println(err)
	}
	defer archivo.Close()
	archivo.Seek(int64(index), 0)
	err = binary.Read(archivo, binary.LittleEndian, &SuperBloqueActual.s_filesystem_type)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &SuperBloqueActual.s_inodes_count)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &SuperBloqueActual.s_blocks_count)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &SuperBloqueActual.s_free_blocks_count)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &SuperBloqueActual.s_free_inodes_count)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &SuperBloqueActual.s_mtime)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &SuperBloqueActual.s_mnt_count)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &SuperBloqueActual.s_magic)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &SuperBloqueActual.s_inode_size)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &SuperBloqueActual.s_block_size)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &SuperBloqueActual.s_first_ino)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &SuperBloqueActual.s_first_blo)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &SuperBloqueActual.s_bm_inode_start)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &SuperBloqueActual.s_bm_block_start)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &SuperBloqueActual.s_inode_start)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &SuperBloqueActual.s_block_start)
	if err != nil {
		fmt.Println(err)
	}
}

func InitInodo(NewInodo *TablaInodo) {
	var err error
	NewInodo.i_uid = make([]byte, 4)
	NewInodo.i_gid = make([]byte, 4)
	NewInodo.i_size = make([]byte, 4)
	NewInodo.i_atime, err = (time.Now()).MarshalBinary()
	if err != nil {
		fmt.Println(err)
	}
	NewInodo.i_ctime, err = (time.Now()).MarshalBinary()
	if err != nil {
		fmt.Println(err)
	}
	NewInodo.i_mtime, err = (time.Now()).MarshalBinary()
	if err != nil {
		fmt.Println(err)
	}
	var newBlocks [16]int
	NewInodo.i_block = Int16ArrayToByteArray(newBlocks)
	NewInodo.i_type = make([]byte, 1)
	NewInodo.i_perm = make([]byte, 4)
}

func EscribirInodo(Path string, InodoActual TablaInodo, index int) {
	archivo, err := os.OpenFile(Path, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println(err)
	}
	defer archivo.Close()
	archivo.Seek(int64(index), 0)
	err = binary.Write(archivo, binary.LittleEndian, &InodoActual.i_uid)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &InodoActual.i_gid)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &InodoActual.i_size)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &InodoActual.i_atime)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &InodoActual.i_ctime)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &InodoActual.i_mtime)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &InodoActual.i_block)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &InodoActual.i_type)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(archivo, binary.LittleEndian, &InodoActual.i_perm)
	if err != nil {
		fmt.Println(err)
	}
}

func LeerInodo(Path string, InodoActual *TablaInodo, index int) {
	archivo, err := os.Open(Path)
	if err != nil {
		fmt.Println(err)
	}
	defer archivo.Close()
	archivo.Seek(int64(index), 0)
	err = binary.Read(archivo, binary.LittleEndian, &InodoActual.i_uid)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &InodoActual.i_gid)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &InodoActual.i_size)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &InodoActual.i_atime)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &InodoActual.i_ctime)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &InodoActual.i_mtime)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &InodoActual.i_block)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &InodoActual.i_type)
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Read(archivo, binary.LittleEndian, &InodoActual.i_perm)
	if err != nil {
		fmt.Println(err)
	}
}

func InitBloqueArchivos(NewBlock *BloqueArchivos) {
	NewBlock.b_content = make([]byte, 64)
}

func EscribirBloqueArchivos(Path string, BlockActual BloqueArchivos, index int) {
	archivo, err := os.OpenFile(Path, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println(err)
	}
	defer archivo.Close()
	archivo.Seek(int64(index), 0)
	err = binary.Write(archivo, binary.LittleEndian, &BlockActual.b_content)
	if err != nil {
		fmt.Println(err)
	}
}

func LeerBloqueArchivos(Path string, BlockActual *BloqueArchivos, index int) {
	archivo, err := os.Open(Path)
	if err != nil {
		fmt.Println(err)
	}
	defer archivo.Close()
	archivo.Seek(int64(index), 0)
	err = binary.Read(archivo, binary.LittleEndian, &BlockActual.b_content)
	if err != nil {
		fmt.Println(err)
	}
}

func InitBloqueCarpeta(NewBlock *BloqueCarpeta) {
	for i := 0; i < 4; i++ {
		NewBlock.b_content[i].b_name = make([]byte, 12)
		NewBlock.b_content[i].b_inodo = make([]byte, 4)
	}
}

func EscribirBloqueCarpeta(Path string, BlockActual BloqueCarpeta, index int) {
	archivo, err := os.OpenFile(Path, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println(err)
	}
	defer archivo.Close()
	archivo.Seek(int64(index), 0)
	for i := 0; i < 4; i++ {
		err = binary.Write(archivo, binary.LittleEndian, &BlockActual.b_content[i].b_name)
		if err != nil {
			fmt.Println(err)
		}
		err = binary.Write(archivo, binary.LittleEndian, &BlockActual.b_content[i].b_inodo)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func LeerBloqueCarpeta(Path string, BlockActual *BloqueCarpeta, index int) {
	archivo, err := os.Open(Path)
	if err != nil {
		fmt.Println(err)
	}
	defer archivo.Close()
	archivo.Seek(int64(index), 0)
	for i := 0; i < 4; i++ {
		err = binary.Read(archivo, binary.LittleEndian, &BlockActual.b_content[i].b_name)
		if err != nil {
			fmt.Println(err)
		}
		err = binary.Read(archivo, binary.LittleEndian, &BlockActual.b_content[i].b_inodo)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func EscribirBM(Path string, index int, Value byte) {
	b := make([]byte, 1)
	b[0] = Value
	archivo, err := os.OpenFile(Path, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println(err)
	}
	defer archivo.Close()
	archivo.Seek(int64(index), 0)
	err = binary.Write(archivo, binary.LittleEndian, &b)
	if err != nil {
		fmt.Println(err)
	}
}

func LeerBM(Path string, index int) byte {
	archivo, err := os.Open(Path)
	if err != nil {
		fmt.Println(err)
	}
	defer archivo.Close()
	b := make([]byte, 1)
	archivo.Seek(int64(index), 0)
	err = binary.Read(archivo, binary.LittleEndian, &b)
	if err != nil {
		fmt.Println(err)
	}
	return b[0]
}

func Int16ArrayToByteArray(intArr [16]int) []byte {
	byteArr := make([]byte, 64)
	for i, num := range intArr {
		binary.LittleEndian.PutUint32(byteArr[i*4:], uint32(num))
	}
	return byteArr
}

func ByteArrayToInt16Array(byteArr []byte) [16]int {
	var intArr [16]int
	for i := 0; i < 16; i++ {
		intArr[i] = int(binary.LittleEndian.Uint32(byteArr[i*4 : (i+1)*4]))
	}
	return intArr
}

//Clone

func (p *Partition) clone() Partition {
	return Partition{
		part_status: append([]byte{}, p.part_status...),
		part_type:   append([]byte{}, p.part_type...),
		part_fit:    append([]byte{}, p.part_fit...),
		part_start:  append([]byte{}, p.part_start...),
		part_size:   append([]byte{}, p.part_size...),
		part_name:   append([]byte{}, p.part_name...),
	}
}

func (m *MBR) clone() MBR {
	return MBR{
		mbr_tamano:         append([]byte{}, m.mbr_tamano...),
		mbr_fecha_creacion: append([]byte{}, m.mbr_fecha_creacion...),
		mbr_dsk_signature:  append([]byte{}, m.mbr_dsk_signature...),
		mbr_dsk_fit:        append([]byte{}, m.mbr_dsk_fit...),
		mbr_partition_1:    m.mbr_partition_1.clone(),
		mbr_partition_2:    m.mbr_partition_2.clone(),
		mbr_partition_3:    m.mbr_partition_3.clone(),
		mbr_partition_4:    m.mbr_partition_4.clone(),
	}
}

func (p *EBR) clone() EBR {
	return EBR{
		part_status: append([]byte{}, p.part_status...),
		part_fit:    append([]byte{}, p.part_fit...),
		part_start:  append([]byte{}, p.part_start...),
		part_size:   append([]byte{}, p.part_size...),
		part_next:   append([]byte{}, p.part_next...),
		part_name:   append([]byte{}, p.part_name...),
	}
}

func (s *SuperBloque) clone() SuperBloque {
	return SuperBloque{
		s_filesystem_type:   append([]byte{}, s.s_filesystem_type...),
		s_inodes_count:      append([]byte{}, s.s_inodes_count...),
		s_blocks_count:      append([]byte{}, s.s_blocks_count...),
		s_free_blocks_count: append([]byte{}, s.s_free_blocks_count...),
		s_free_inodes_count: append([]byte{}, s.s_free_inodes_count...),
		s_mtime:             append([]byte{}, s.s_mtime...),
		s_mnt_count:         append([]byte{}, s.s_mnt_count...),
		s_magic:             append([]byte{}, s.s_magic...),
		s_inode_size:        append([]byte{}, s.s_inode_size...),
		s_block_size:        append([]byte{}, s.s_block_size...),
		s_first_ino:         append([]byte{}, s.s_first_ino...),
		s_first_blo:         append([]byte{}, s.s_first_blo...),
		s_bm_inode_start:    append([]byte{}, s.s_bm_inode_start...),
		s_bm_block_start:    append([]byte{}, s.s_bm_block_start...),
		s_inode_start:       append([]byte{}, s.s_inode_start...),
		s_block_start:       append([]byte{}, s.s_block_start...),
	}
}

func (t *TablaInodo) clone() TablaInodo {
	return TablaInodo{
		i_uid:   append([]byte{}, t.i_uid...),
		i_gid:   append([]byte{}, t.i_gid...),
		i_size:  append([]byte{}, t.i_size...),
		i_atime: append([]byte{}, t.i_atime...),
		i_ctime: append([]byte{}, t.i_ctime...),
		i_mtime: append([]byte{}, t.i_mtime...),
		i_block: append([]byte{}, t.i_block...),
		i_type:  append([]byte{}, t.i_type...),
		i_perm:  append([]byte{}, t.i_perm...),
	}
}

func (c *content) clone() content {
	return content{
		b_name:  append([]byte{}, c.b_name...),
		b_inodo: append([]byte{}, c.b_inodo...),
	}
}

func (b *BloqueCarpeta) clone() BloqueCarpeta {
	var bc BloqueCarpeta
	for i := range b.b_content {
		bc.b_content[i] = b.b_content[i].clone()
	}
	return bc
}

func (b *BloqueArchivos) clone() BloqueArchivos {
	return BloqueArchivos{
		b_content: append([]byte{}, b.b_content...),
	}
}

func (mp *MountPart) clone() MountPart {
	return MountPart{
		Active:     mp.Active,
		ID:         mp.ID,
		Path:       mp.Path,
		part_name:  append([]byte{}, mp.part_name...),
		part_size:  append([]byte{}, mp.part_size...),
		part_start: append([]byte{}, mp.part_start...),
	}
}

func (s *Session) clone() Session {
	return Session{
		User:   s.User,
		IDU:    s.IDU,
		Grupo:  s.Grupo,
		IDG:    s.IDG,
		Active: s.Active.clone(),
	}
}

//See

func VerDisco(Path string) {
	var MBRDsk MBR
	Cabeza := -1
	InitMBR(&MBRDsk)
	LeerMBR(Path, &MBRDsk)
	Reporte := ""
	Reporte += "Size: " + strconv.Itoa(int(binary.LittleEndian.Uint32(MBRDsk.mbr_tamano))) + "\n"
	var Tiempo time.Time
	Tiempo.UnmarshalBinary(MBRDsk.mbr_fecha_creacion)
	Reporte += "Fecha de creación: " + Tiempo.String() + "\n"
	Reporte += "Signature: " + strconv.Itoa(int(binary.LittleEndian.Uint32(MBRDsk.mbr_dsk_signature))) + "\n"
	Reporte += "Fit: " + string(MBRDsk.mbr_dsk_fit) + "\n"
	if MBRDsk.mbr_partition_1.part_status[0] == 'A' {
		Reporte += "\tParticion 1." + "\n"
		Reporte += "\tType: " + string(MBRDsk.mbr_partition_1.part_type) + "\n"
		Reporte += "\tFit: " + string(MBRDsk.mbr_partition_1.part_fit) + "\n"
		Reporte += "\tStart: " + strconv.Itoa(int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_1.part_start))) + "\n"
		Reporte += "\tSize: " + strconv.Itoa(int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_1.part_size))) + "\n"
		Reporte += "\tName: " + string(MBRDsk.mbr_partition_1.part_name) + "\n"
		if MBRDsk.mbr_partition_1.part_type[0] == 'E' {
			Cabeza = int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_1.part_start))
		}
		if MBRDsk.mbr_partition_2.part_status[0] == 'A' {
			Reporte += "\tParticion 2." + "\n"
			Reporte += "\tType: " + string(MBRDsk.mbr_partition_2.part_type) + "\n"
			Reporte += "\tFit: " + string(MBRDsk.mbr_partition_2.part_fit) + "\n"
			Reporte += "\tStart: " + strconv.Itoa(int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_2.part_start))) + "\n"
			Reporte += "\tSize: " + strconv.Itoa(int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_2.part_size))) + "\n"
			Reporte += "\tName: " + string(MBRDsk.mbr_partition_2.part_name) + "\n"
			if MBRDsk.mbr_partition_2.part_type[0] == 'E' {
				Cabeza = int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_2.part_start))
			}
			if MBRDsk.mbr_partition_3.part_status[0] == 'A' {
				Reporte += "\tParticion 3." + "\n"
				Reporte += "\tType: " + string(MBRDsk.mbr_partition_3.part_type) + "\n"
				Reporte += "\tFit: " + string(MBRDsk.mbr_partition_3.part_fit) + "\n"
				Reporte += "\tStart: " + strconv.Itoa(int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_3.part_start))) + "\n"
				Reporte += "\tSize: " + strconv.Itoa(int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_3.part_size))) + "\n"
				Reporte += "\tName: " + string(MBRDsk.mbr_partition_3.part_name) + "\n"
				if MBRDsk.mbr_partition_3.part_type[0] == 'E' {
					Cabeza = int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_3.part_start))
				}
				if MBRDsk.mbr_partition_4.part_status[0] == 'A' {
					Reporte += "\tParticion 4." + "\n"
					Reporte += "\tType: " + string(MBRDsk.mbr_partition_4.part_type) + "\n"
					Reporte += "\tFit: " + string(MBRDsk.mbr_partition_4.part_fit) + "\n"
					Reporte += "\tStart: " + strconv.Itoa(int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_4.part_start))) + "\n"
					Reporte += "\tSize: " + strconv.Itoa(int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_4.part_size))) + "\n"
					Reporte += "\tName: " + string(MBRDsk.mbr_partition_4.part_name) + "\n"
					if MBRDsk.mbr_partition_4.part_type[0] == 'E' {
						Cabeza = int(binary.LittleEndian.Uint32(MBRDsk.mbr_partition_4.part_start))
					}
				}
			}
		}
	}
	if Cabeza != -1 {
		Aux := Cabeza
		Cont := 1
		for Aux != 0 {
			var EBRActual EBR
			InitEBR(&EBRActual)
			LeerEBR(Path, &EBRActual, Aux)
			Reporte += "\t\tParticion Logica " + strconv.Itoa(Cont) + "\n"
			Reporte += "\t\tEstatus " + string(EBRActual.part_status) + "\n"
			Reporte += "\t\tFit " + string(EBRActual.part_fit) + "\n"
			Reporte += "\t\tStart " + strconv.Itoa(int(binary.LittleEndian.Uint32(EBRActual.part_start))) + "\n"
			Reporte += "\t\tSize " + strconv.Itoa(int(binary.LittleEndian.Uint32(EBRActual.part_size))) + "\n"
			Reporte += "\t\tNext " + strconv.Itoa(int(binary.LittleEndian.Uint32(EBRActual.part_next))) + "\n"
			Reporte += "\t\tFit " + string(EBRActual.part_name) + "\n"
			Aux = int(binary.LittleEndian.Uint32(EBRActual.part_next))
			Cont++
		}
	}
	Console += Reporte
}

func VerMounts() {
	for i := 0; i < 10; i++ {
		if ActivePart[i].Active {
			Console += "Partition No. " + strconv.Itoa(i+1) + ":\n"
			Console += "\tID: " + ActivePart[i].ID + "\n"
			Console += "\tPath: " + ActivePart[i].Path + "\n"
			Console += "\tName: " + string(ActivePart[i].part_name) + "\n"
			Console += "\tSize: " + strconv.Itoa(int(binary.LittleEndian.Uint32(ActivePart[i].part_size))) + "\n"
			Console += "\tStart: " + strconv.Itoa(int(binary.LittleEndian.Uint32(ActivePart[i].part_start))) + "\n"
		}
	}
}

func VerInfo(ID string) {
	ActiveParticion := -1
	for i := 0; i < 10; i++ {
		if ActivePart[i].Active && ID == ActivePart[i].ID {
			ActiveParticion = i
			break
		}
	}
	if ActiveParticion != -1 {
		Activa := ActivePart[ActiveParticion]
		//PartitionSize := int(binary.LittleEndian.Uint32(Activa.part_size))
		PartitionStart := int(binary.LittleEndian.Uint32(Activa.part_start))
		var SB SuperBloque
		InitSuperBloque(&SB)
		LeerSuperBloque(Activa.Path, &SB, PartitionStart)
		InodesCount := int(binary.LittleEndian.Uint32(SB.s_inodes_count))
		BlocksCount := int(binary.LittleEndian.Uint32(SB.s_blocks_count))
		InodeStart := int(binary.LittleEndian.Uint32(SB.s_inode_start))
		BlockStart := int(binary.LittleEndian.Uint32(SB.s_block_start))
		BMInodeStart := int(binary.LittleEndian.Uint32(SB.s_bm_inode_start))
		BMBlockStart := int(binary.LittleEndian.Uint32(SB.s_bm_block_start))
		Console += "SuperBloque: \n"
		Console += "\tFileSystem Type: " + strconv.Itoa(int(binary.LittleEndian.Uint32(SB.s_filesystem_type))) + "\n"
		Console += "\tInodes Count: " + strconv.Itoa(InodesCount) + "\n"
		Console += "\tBlocks Count: " + strconv.Itoa(BlocksCount) + "\n"
		Console += "\tFree Blocks Count: " + strconv.Itoa(int(binary.LittleEndian.Uint32(SB.s_free_blocks_count))) + "\n"
		Console += "\tFree Inodes Count: " + strconv.Itoa(int(binary.LittleEndian.Uint32(SB.s_free_inodes_count))) + "\n"
		var Tiempo time.Time
		Tiempo.UnmarshalBinary(SB.s_mtime)
		Console += "\tFecha de creación: " + Tiempo.String() + "\n"
		Console += "\tmnt_count: " + strconv.Itoa(int(binary.LittleEndian.Uint32(SB.s_mnt_count))) + "\n"
		Console += "\tMagic: " + strconv.Itoa(int(binary.LittleEndian.Uint32(SB.s_magic))) + "\n"
		Console += "\tInode_size: " + strconv.Itoa(int(binary.LittleEndian.Uint32(SB.s_inode_size))) + "\n"
		Console += "\tBlock_size: " + strconv.Itoa(int(binary.LittleEndian.Uint32(SB.s_block_size))) + "\n"
		Console += "\tFirst_ino: " + strconv.Itoa(int(binary.LittleEndian.Uint32(SB.s_first_ino))) + "\n"
		Console += "\tFirst_blo: " + strconv.Itoa(int(binary.LittleEndian.Uint32(SB.s_first_blo))) + "\n"
		Console += "\tBM_inode_start: " + strconv.Itoa(BMInodeStart) + "\n"
		Console += "\tBM_block_start: " + strconv.Itoa(BMBlockStart) + "\n"
		Console += "\tinode_start: " + strconv.Itoa(InodeStart) + "\n"
		Console += "\tblock_start: " + strconv.Itoa(BlockStart) + "\n"
		Console += "\tBitMap Inodos: \n\t\t"
		for i := 0; i < InodesCount; i++ {
			if LeerBM(Activa.Path, BMInodeStart+i) == 1 {
				Console += "1"
			} else {
				Console += "0"
			}
			if i%25 != 24 && i != InodesCount-1 {
				Console += ","
			} else if i != InodesCount-1 {
				Console += "\n\t\t"
			} else {
				Console += "\n"
			}
		}
		Console += "\tBitMap Blocks: \n\t\t"
		for i := 0; i < BlocksCount; i++ {
			if LeerBM(Activa.Path, BMBlockStart+i) == 1 {
				Console += "1"
			} else {
				Console += "0"
			}
			if i%25 != 24 && i != BlocksCount-1 {
				Console += ","
			} else if i != BlocksCount-1 {
				Console += "\n\t\t"
			} else {
				Console += "\n"
			}
		}
		Console += "\tInodos: \n"
		for i := 0; i < InodesCount; i++ {
			if LeerBM(Activa.Path, BMInodeStart+i) == 1 {
				var InodoActual TablaInodo
				InitInodo(&InodoActual)
				LeerInodo(Activa.Path, &InodoActual, InodeStart+i*int(unsafe.Sizeof(InodoActual)))
				Console += "\t\tInodo " + strconv.Itoa(i) + "\n"
				Console += "\t\t\tUID: " + strconv.Itoa(int(binary.LittleEndian.Uint32(InodoActual.i_uid))) + "\n"
				Console += "\t\t\tGID: " + strconv.Itoa(int(binary.LittleEndian.Uint32(InodoActual.i_gid))) + "\n"
				Console += "\t\t\tSize: " + strconv.Itoa(int(binary.LittleEndian.Uint32(InodoActual.i_size))) + "\n"
				Tiempo.UnmarshalBinary(InodoActual.i_atime)
				Console += "\t\t\taTime: " + Tiempo.String() + "\n"
				Tiempo.UnmarshalBinary(InodoActual.i_ctime)
				Console += "\t\t\tcTime: " + Tiempo.String() + "\n"
				Tiempo.UnmarshalBinary(InodoActual.i_mtime)
				Console += "\t\t\tmTime: " + Tiempo.String() + "\n"
				TempBlocks := ByteArrayToInt16Array(InodoActual.i_block)
				for j := 0; j < 16; j++ {
					if TempBlocks[j] != 4294967295 {
						Console += "\t\t\tBlock " + strconv.Itoa(j+1) + ": " + strconv.Itoa(TempBlocks[j]) + "\n"
					} else {
						Console += "\t\t\tBlock " + strconv.Itoa(j+1) + ": -1\n"
					}
				}
				if InodoActual.i_type[0] == 0 {
					Console += "\t\t\tType: 0\n"
				} else {
					Console += "\t\t\tType: 1\n"
				}
				Console += "\t\t\tPerm: " + strconv.Itoa(int(binary.LittleEndian.Uint32(InodoActual.i_perm))) + "\n"
			}
		}
		Console += "\tBloques: \n"
		for i := 0; i < BlocksCount; i++ {
			if LeerBM(Activa.Path, BMBlockStart+i) == 1 {
				var BA BloqueArchivos
				InitBloqueArchivos(&BA)
				LeerBloqueArchivos(Activa.Path, &BA, BlockStart+i*64)
				var BC BloqueCarpeta
				InitBloqueCarpeta(&BC)
				LeerBloqueCarpeta(Activa.Path, &BC, BlockStart+i*64)
				Console += "\t\tBloque " + strconv.Itoa(i+1) + "\n"
				Console += "\t\t\tArchivo_Content: " + string(BA.b_content) + "\n"
				for j := 0; j < 4; j++ {
					Console += "\t\tCarpeta_" + strconv.Itoa(j+1) + "_name: " + string(BC.b_content[j].b_name) + "\n"
					if int(binary.LittleEndian.Uint32(BC.b_content[j].b_inodo)) != 4294967295 {
						Console += "\t\tCarpeta_" + strconv.Itoa(j+1) + "_inodo: " + strconv.Itoa(int(binary.LittleEndian.Uint32(BC.b_content[j].b_inodo))) + "\n"
					} else {
						Console += "\t\tCarpeta_" + strconv.Itoa(j+1) + "_inodo: -1\n"
					}
				}
			}
		}
		Temp := Sesion.Active.clone()
		Sesion.Active = Activa.clone()
		var InodoActual TablaInodo
		InitInodo(&InodoActual)
		Usuarios_txt := LeerArchivo(InodeStart + int(unsafe.Sizeof(InodoActual)))
		if Usuarios_txt != "" {
			Console += "\tUsuarios: \n"
			Console += Usuarios_txt
		}
		Sesion.Active = Temp.clone()
	} else {
		Console += "Error, no se encontro la ID\n"
	}
}
