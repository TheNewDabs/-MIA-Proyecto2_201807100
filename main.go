package main

import (
	"encoding/binary"
	"fmt"
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

type SB struct {
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

const BM0 = 0
const BM1 = 1

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
seedisk >path=/home/dabs/201807100/Disco1.dk

mkdisk >size=25 >fit=bf >unit=m >path="/home/dabs/201807100/primer semestre/Disco2.dk"
fdisk >size=500 >unit=k >path="/home/dabs/201807100/primer semestre/Disco2.dk" >name=Particion1 >fit=ff
fdisk >size=1024 >path="/home/dabs/201807100/primer semestre/Disco2.dk" >unit=k >name=Particion2
fdisk >size=10 >unit=m >path="/home/dabs/201807100/primer semestre/Disco2.dk" >name=Particion3
fdisk >unit=k >size=4096 >path="/home/dabs/201807100/primer semestre/Disco2.dk" >type=E >name=Particion4 >fit=wf
seedisk >path="/home/dabs/201807100/primer semestre/Disco2.dk"

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
seemounts`)
	fmt.Println("Consola: ")
	fmt.Println(Console[:len(Console)-1])
}

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
		} else {
			Console += "Error, comando desconocido\n"
		}
	}
	return true
}

func CrearDisco(Size int, Path string, Fit string, Unit string) {
	var err error
	var NewMBR MBR
	InitMBR(&NewMBR)
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
		fmt.Println(strings.Compare(NameA, Name))
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

func getName(Path string) string {
	for strings.Index(Path, "/") != -1 {
		Path = Path[strings.Index(Path, "/")+1:]
	}
	if strings.Index(Path, ".") != -1 {
		Path = Path[:strings.Index(Path, ".")]
	}
	return Path
}

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

func Init16String(Original string) string {
	Temp := make([]byte, 16)
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

func InitMountPart(NewMountPart *MountPart) {
	NewMountPart.part_start = make([]byte, 4)
	NewMountPart.part_size = make([]byte, 4)
	NewMountPart.part_name = make([]byte, 16)
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

func (s *SB) clone() SB {
	return SB{
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
