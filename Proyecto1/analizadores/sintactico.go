package analizadores

import (
	"Proyecto1/estructuras"
	"Proyecto1/funciones"
	"fmt"

	"github.com/doun/terminal/color"
)

var (
	syntaxError                                                                          bool = false
	token                                                                                int  = -1
	tokenAux                                                                             *estructuras.Token
	vSize, vPath, vName, vUnit, vType, vFit, vDelete, vAdd, vID, vNombre, vRuta, vFormat string = "", "", "", "", "", "", "", "", "", "", "", ""
	ejMkdisk, ejFdisk, ejRmdisk, ejMount, ejUnmount, ejReporte, ejMkfs                   bool   = false, false, false, false, false, false, false
	//ListaIDs para desmontar IDs
	ListaIDs []string
	//Discos para almacenar discos son particiones montadas

)

func resetearBanderas() {
	ejFdisk = false
	ejMkdisk = false
	ejRmdisk = false
	ejMount = false
	ejUnmount = false
	ejReporte = false
	ejMkfs = false
}

func resetearValores() {
	vSize = ""
	vPath = ""
	vName = ""
	vUnit = ""
	vType = ""
	vFit = ""
	vDelete = ""
	vAdd = ""
	vID = ""
	vNombre = ""
	vRuta = ""
	vFormat = ""
}

//Sintactico fuction
func Sintactico() {
	syntaxError = false
	tokenAux = nextToken()
	token = -1

	if token < (len(tokens) - 1) {
		tokenAux = nextToken()
		inicio()
	}

	if !syntaxError && token >= (len(tokens)-1) {
		color.Println("@{!c}	Analisis sintáctico exitoso")
	} else {
		color.Println("@{!r}	Error sintáctico encontrado")
	}

}

func inicio() {

	if tokenAux.GetTipo() == "TK_CMT" {
		color.Printf("@{!m}%v\n", tokenAux.GetLexema())
		otraInstruccion()
	} else if tokenAux.GetTipo() == "TK_PAUSE" {
		Pausa()
		otraInstruccion()
	} else if tokenAux.GetTipo() == "TK_EXEC" {

		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_PATH") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ASIG") {
				tokenAux = nextToken()
				if tokenCorrecto(tokenAux, "TK_FILE") {
					//LEER ARCHIVO
					otraInstruccion()
				} else {
					syntaxError = true
				}
			} else {
				syntaxError = true
			}
		} else {
			syntaxError = true
		}

	} else if tokenAux.GetTipo() == "TK_RMDISK" {

		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_PATH") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ASIG") {
				tokenAux = nextToken()
				if tokenCorrecto(tokenAux, "TK_FILE") {
					//BORRAR DISCO
					vPath = tokenAux.GetLexema()
					funciones.EjecutarRmDisk(vPath)
					resetearBanderas()
					resetearValores()
					otraInstruccion()
				} else {
					syntaxError = true
				}
			} else {
				syntaxError = true
			}
		} else {
			syntaxError = true
		}

	} else if tokenAux.GetTipo() == "TK_MKDISK" {
		ejMkdisk = true
		paramMkDisk()
		if ejMkdisk {
			funciones.EjecutarMkDisk(vSize, vPath, vName, vUnit)
			resetearBanderas()
			resetearValores()
		}
		otraInstruccion()
	} else if tokenAux.GetTipo() == "TK_FDISK" {
		ejFdisk = true
		paramFDisk()
		if ejFdisk {
			funciones.EjecutarFDisk(vSize, vUnit, vPath, vType, vFit, vDelete, vName, vAdd)
			resetearBanderas()
			resetearValores()
		}
		otraInstruccion()
	} else if tokenAux.GetTipo() == "TK_MNT" {
		ejMount = true
		paramMount()
		if ejMount {
			funciones.EjecutarMount(vPath, vName)
			resetearBanderas()
			resetearValores()
		}
		otraInstruccion()
	} else if tokenAux.GetTipo() == "TK_REP" {
		ejReporte = true
		paramRep()
		if ejReporte {
			funciones.EjecutarReporte(vNombre, vPath, vRuta, vID)
			resetearBanderas()
			resetearValores()
		}
		otraInstruccion()
	} else if tokenAux.GetTipo() == "TK_MKFS" {
		ejMkfs = true
		paramMkfs()
		if ejMkfs {
			funciones.EjecutarMkfs()
			resetearBanderas()
			resetearValores()
		}
		otraInstruccion()
	} else if tokenAux.GetTipo() == "TK_UMNT" {
		ListaIDs = nil
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_PID") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_NUM") {
				tokenAux = nextToken()
				if tokenCorrecto(tokenAux, "TK_ASIG") {
					tokenAux = nextToken()
					if tokenCorrecto(tokenAux, "TK_ID") {
						ejUnmount = true
						//Guardar ID
						ListaIDs = append(ListaIDs, tokenAux.GetLexema())
						otroID()
						if ejUnmount {
							funciones.EjecutarUnmount(&ListaIDs)
							resetearBanderas()
							resetearValores()
						}
						otraInstruccion()
					} else {
						syntaxError = true
					}
				} else {
					syntaxError = true
				}
			} else {
				syntaxError = true
			}
		} else {
			syntaxError = true
		}

	} else {
		syntaxError = true
		fmt.Println("Se esperaba fdisk, mkdisk, mount, etc.")
	}

}

func paramMkDisk() {
	tokenAux = nextToken()
	if tokenCorrecto(tokenAux, "TK_SIZE") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_NUM") {
				//SETEAR SIZE
				vSize = tokenAux.GetLexema()
				otroParamMkDisk()

			} else {
				ejMkdisk = false
				syntaxError = true
			}
		} else {
			ejMkdisk = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_PATH") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_DIR") {
				//SETEAR PATH
				vPath = tokenAux.GetLexema()
				otroParamMkDisk()
			} else {
				ejMkdisk = false
				syntaxError = true
			}
		} else {
			ejMkdisk = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_NAME") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR NAME
				vName = tokenAux.GetLexema()
				otroParamMkDisk()
			} else {
				ejMkdisk = false
				syntaxError = true
			}
		} else {
			ejMkdisk = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_UNIT") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_BYTES") {
				//SETEAR BYTES
				vUnit = tokenAux.GetLexema()
				otroParamMkDisk()
			} else {
				ejMkdisk = false
				syntaxError = true
			}
		} else {
			ejMkdisk = false
			syntaxError = true
		}
	} else {
		ejMkdisk = false
		syntaxError = true
		fmt.Println("Se esperaba -size, -path, -name, etc.")
	}
}

func otroParamMkDisk() {
	if token < (len(tokens) - 1) {
		if tokens[token+1].GetTipo() == "TK_SIZE" || tokens[token+1].GetTipo() == "TK_PATH" || tokens[token+1].GetTipo() == "TK_NAME" || tokens[token+1].GetTipo() == "TK_UNIT" {
			paramMkDisk()
		}
	}
}

func paramFDisk() {
	tokenAux = nextToken()
	if tokenCorrecto(tokenAux, "TK_SIZE") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_NUM") {
				//SETEAR SIZE
				vSize = tokenAux.GetLexema()
				otroParamFDisk()
			} else {
				ejFdisk = false
				syntaxError = true
			}
		} else {
			ejFdisk = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_UNIT") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_BYTES") {
				//SETEAR BYTES
				vUnit = tokenAux.GetLexema()
				otroParamFDisk()
			} else {
				ejFdisk = false
				syntaxError = true
			}
		} else {
			ejFdisk = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_PATH") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_FILE") {
				//SETEAR PATH
				vPath = tokenAux.GetLexema()
				otroParamFDisk()
			} else {
				ejFdisk = false
				syntaxError = true
			}
		} else {
			ejFdisk = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_TYPE") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_PEL") {
				//SETEAR TYPE
				vType = tokenAux.GetLexema()
				otroParamFDisk()
			} else {
				ejFdisk = false
				syntaxError = true
			}
		} else {
			ejFdisk = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_FIT") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_BFW") {
				//SETEAR FIT
				vFit = tokenAux.GetLexema()
				otroParamFDisk()
			} else {
				ejFdisk = false
				syntaxError = true
			}
		} else {
			ejFdisk = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_DEL") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_FF") {
				//SETEAR DELETE MODE
				vDelete = tokenAux.GetLexema()
				otroParamFDisk()
			} else {
				ejFdisk = false
				syntaxError = true
			}
		} else {
			ejFdisk = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_NAME") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR NAME
				vName = tokenAux.GetLexema()
				otroParamFDisk()
			} else {
				ejFdisk = false
				syntaxError = true
			}
		} else {
			ejFdisk = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_ADD") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_NUM") {
				//SETEAR NUM
				vAdd = tokenAux.GetLexema()
				otroParamFDisk()
			} else {
				ejFdisk = false
				syntaxError = true
			}
		} else {
			ejFdisk = false
			syntaxError = true
		}
	} else {
		ejFdisk = false
		syntaxError = true
		fmt.Println("Se esperaba -size, -path, -name, etc.")
	}
}

func otroParamFDisk() {
	if token < (len(tokens) - 1) {
		if tokens[token+1].GetTipo() == "TK_SIZE" || tokens[token+1].GetTipo() == "TK_PATH" || tokens[token+1].GetTipo() == "TK_NAME" || tokens[token+1].GetTipo() == "TK_UNIT" || tokens[token+1].GetTipo() == "TK_TYPE" || tokens[token+1].GetTipo() == "TK_FIT" || tokens[token+1].GetTipo() == "TK_DEL" || tokens[token+1].GetTipo() == "TK_ADD" {
			paramFDisk()
		}
	}
}

func paramMount() {
	tokenAux = nextToken()
	if tokenCorrecto(tokenAux, "TK_PATH") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_FILE") {
				//SETEAR PATH
				vPath = tokenAux.GetLexema()
				otroParamMount()
			} else {
				ejMount = false
				syntaxError = true
			}
		} else {
			ejMount = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_NAME") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR NAME
				vName = tokenAux.GetLexema()
				otroParamMount()
			} else {
				ejMount = false
				syntaxError = true
			}
		} else {
			ejMount = false
			syntaxError = true
		}
	} else {
		ejMount = false
		syntaxError = true
		fmt.Println("Se esperaba -path o -name.")
	}
}

func otroParamMount() {
	if token < (len(tokens) - 1) {
		if tokens[token+1].GetTipo() == "TK_PATH" || tokens[token+1].GetTipo() == "TK_NAME" {
			paramMount()
		}
	}
}

func paramRep() {
	tokenAux = nextToken()
	if tokenCorrecto(tokenAux, "TK_NOMBRE") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_TIPOREPORTE") {
				//SETEAR TIPOREPORTE
				vNombre = tokenAux.GetLexema()
				otroParamRep()
			} else {
				ejReporte = false
				syntaxError = true
			}
		} else {
			ejReporte = false
			syntaxError = true
		}

	} else if tokenCorrecto(tokenAux, "TK_PATH") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_FILE") {
				//SETEAR PATH
				vPath = tokenAux.GetLexema()
				otroParamRep()
			} else {
				ejReporte = false
				syntaxError = true
			}
		} else {
			ejReporte = false
			syntaxError = true
		}

	} else if tokenCorrecto(tokenAux, "TK_PID") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR ID
				vID = tokenAux.GetLexema()
				otroParamRep()
			} else {
				ejReporte = false
				syntaxError = true
			}
		} else {
			ejReporte = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_RUTA") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_FILE") || tokenCorrecto(tokenAux, "TK_DIR") {
				//SETEAR RUTA
				vRuta = tokenAux.GetLexema()
				otroParamRep()
			} else {
				ejReporte = false
				syntaxError = true
			}
		} else {
			ejReporte = false
			syntaxError = true
		}
	} else {
		ejReporte = false
		syntaxError = true
		fmt.Println("Se esperaba -nombre, -path, -id, etc.")
	}
}

func otroParamRep() {
	if token < (len(tokens) - 1) {
		if tokens[token+1].GetTipo() == "TK_NOMBRE" || tokens[token+1].GetTipo() == "TK_PATH" || tokens[token+1].GetTipo() == "TK_PID" || tokens[token+1].GetTipo() == "TK_RUTA" {
			paramRep()
		}
	}
}

func otroID() {
	if token < (len(tokens) - 1) {
		if tokens[token+1].GetTipo() == "TK_PID" {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_PID") {
				tokenAux = nextToken()
				if tokenCorrecto(tokenAux, "TK_NUM") {
					tokenAux = nextToken()
					if tokenCorrecto(tokenAux, "TK_ASIG") {
						tokenAux = nextToken()
						if tokenCorrecto(tokenAux, "TK_ID") {
							//Guardar ID
							ListaIDs = append(ListaIDs, tokenAux.GetLexema())
							otroID()
						} else {
							ejUnmount = false
							syntaxError = true
						}
					} else {
						ejUnmount = false
						syntaxError = true
					}
				} else {
					ejUnmount = false
					syntaxError = true
				}
			} else {
				ejUnmount = false
				syntaxError = true
			}
		}
	}
}

func paramMkfs() {
	tokenAux = nextToken()
	if tokenCorrecto(tokenAux, "TK_PID") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR ID
				vID = tokenAux.GetLexema()
				otroParamMkfs()
			} else {
				ejMkfs = false
				syntaxError = true
			}
		} else {
			ejMkfs = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_TYPE") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_FF") {
				//SETEAR FORMAT MODE
				vFormat = tokenAux.GetLexema()
				otroParamMkfs()
			} else {
				ejMkfs = false
				syntaxError = true
			}
		} else {
			ejMkfs = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_ADD") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_NUM") {
				//SETEAR NUM
				vAdd = tokenAux.GetLexema()
				otroParamMkfs()
			} else {
				ejMkfs = false
				syntaxError = true
			}
		} else {
			ejMkfs = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_UNIT") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_BYTES") {
				//SETEAR BYTES
				vUnit = tokenAux.GetLexema()
				otroParamFDisk()
			} else {
				ejMkfs = false
				syntaxError = true
			}
		} else {
			ejMkfs = false
			syntaxError = true
		}
	} else {
		ejMkfs = false
		syntaxError = true
		fmt.Println("Se esperaba -id, -type, etc.")
	}
}

func otroParamMkfs() {
	if token < (len(tokens) - 1) {
		if tokens[token+1].GetTipo() == "TK_PID" || tokens[token+1].GetTipo() == "TK_TYPE" || tokens[token+1].GetTipo() == "TK_ADD" || tokens[token+1].GetTipo() == "TK_UNIT" {
			paramMkfs()
		}
	}
}

func tokenCorrecto(taux *estructuras.Token, tipo string) bool {
	if taux != nil {
		if taux.GetTipo() == tipo {
			return true
		}
		return false
	}
	return false
}

func otraInstruccion() {
	if token < (len(tokens) - 1) {
		tokenAux = nextToken()
		inicio()
	}
}

func nextToken() *estructuras.Token {
	if token < (len(tokens) - 1) {
		token++
		return tokens[token]
	}
	return nil
}

func lastToken() *estructuras.Token {
	if token < (len(tokens) - 1) {
		token--
		return tokens[token]
	}
	return nil
}
