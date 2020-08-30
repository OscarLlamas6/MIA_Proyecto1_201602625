package analizadores

import (
	"Proyecto1/estructuras"
	"Proyecto1/funciones"

	"github.com/doun/terminal/color"
)

var (
	syntaxError                                                                          bool = false
	token                                                                                int  = -1
	tokenAux                                                                             *estructuras.Token
	vSize, vPath, vName, vUnit, vType, vFit, vDelete, vAdd, vID, vNombre, vRuta, vFormat string = "", "", "", "", "", "", "", "", "", "", "", ""
	vUser, vPass, vGroup, vUgo, vR, vCont, vP, vRf                                       string = "", "", "", "", "", "", "", ""
	ejMkdisk, ejFdisk, ejRmdisk, ejMount, ejUnmount, ejReporte, ejMkfs                   bool   = false, false, false, false, false, false, false
	ejLogin, ejMkgrp, ejRmgrp, ejMkusr, ejRmusr, ejChmod, ejMkfile                       bool   = false, false, false, false, false, false, false
	ejCat, ejRm                                                                          bool   = false, false
	//ListaIDs para desmontar IDs
	ListaIDs []string
	//ListaFiles para la function CAT
	ListaFiles []string
)

func resetearBanderas() {
	ejFdisk = false
	ejMkdisk = false
	ejRmdisk = false
	ejMount = false
	ejUnmount = false
	ejReporte = false
	ejMkfs = false
	ejLogin = false
	ejMkgrp = false
	ejRmgrp = false
	ejMkusr = false
	ejRmusr = false
	ejChmod = false
	ejMkfile = false
	ejCat = false
	ejRm = false
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
	vUser = ""
	vPass = ""
	vGroup = ""
	vUgo = ""
	vR = ""
	vCont = ""
	vP = ""
	vRf = ""
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
		//color.Println("@{!c}	Analisis sintáctico exitoso")
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
	} else if tokenAux.GetTipo() == "TK_LOGOUT" {
		funciones.EjecutarLogout()
		otraInstruccion()
	} else if tokenAux.GetTipo() == "TK_EXEC" {

		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_PATH") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ASIG") {
				tokenAux = nextToken()
				if tokenCorrecto(tokenAux, "TK_FILE") {
					//LEER ARCHIVO
					color.Println("@{!r}No se puede ejecutar un script llamado desde otro script.")
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
	} else if tokenAux.GetTipo() == "TK_LOGIN" {
		ejLogin = true
		paramLogin()
		if ejLogin {
			funciones.EjecutarLogin(vUser, vPass, vID)
			resetearBanderas()
			resetearValores()
		}
		otraInstruccion()
	} else if tokenAux.GetTipo() == "TK_MKFS" {
		ejMkfs = true
		paramMkfs()
		if ejMkfs {
			funciones.EjecutarMkfs(vID, vFormat, vAdd, vUnit)
			resetearBanderas()
			resetearValores()
		}
		otraInstruccion()
	} else if tokenAux.GetTipo() == "TK_MKGRP" {
		ejMkgrp = true
		paramMkgrp()
		if ejMkgrp {
			funciones.EjecutarMkgrp(vGroup, vID)
			resetearBanderas()
			resetearValores()
		}
		otraInstruccion()
	} else if tokenAux.GetTipo() == "TK_RMGRP" {
		ejRmgrp = true
		paramRmgrp()
		if ejRmgrp {
			funciones.EjecutarRmgrp(vGroup, vID)
			resetearBanderas()
			resetearValores()
		}
		otraInstruccion()
	} else if tokenAux.GetTipo() == "TK_MKUSR" {
		ejMkusr = true
		paramMkusr()
		if ejMkusr {
			funciones.EjecutarMkusr(vUser, vPass, vGroup, vID)
			resetearBanderas()
			resetearValores()
		}
		otraInstruccion()
	} else if tokenAux.GetTipo() == "TK_RMUSR" {
		ejRmusr = true
		paramRmusr()
		if ejRmusr {
			funciones.EjecutarRmusr(vUser, vID)
			resetearBanderas()
			resetearValores()
		}
		otraInstruccion()
	} else if tokenAux.GetTipo() == "TK_CHMOD" {
		ejChmod = true
		paramChmod()
		if ejChmod {
			funciones.EjecutarChmod(vID, vPath, vUgo, vR)
			resetearBanderas()
			resetearValores()
		}
		otraInstruccion()
	} else if tokenAux.GetTipo() == "TK_MKFILE" {
		ejMkfile = true
		paramMkfile()
		if ejMkfile {
			funciones.EjecutarMkfile(vID, vPath, vSize, vCont, vP)
			resetearBanderas()
			resetearValores()
		}
		otraInstruccion()
	} else if tokenAux.GetTipo() == "TK_CAT" {
		ListaFiles = nil
		ejCat = true
		paramCat()
		if ejCat {
			funciones.EjecutarCat(vID, &ListaFiles)
			resetearBanderas()
			resetearValores()
		}
		otraInstruccion()
	} else if tokenAux.GetTipo() == "TK_RM" {
		ejRm = true
		paramRm()
		if ejRm {
			funciones.EjecutarRm(vID, vPath, vRf)
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
		//color.Println("@{!r}Se esperaba fdisk, mkdisk, mount, etc.")
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
		color.Println("@{!r}Se esperaba -size, -path, -name, etc.")
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
		color.Println("@{!r}Se esperaba -size, -path, -name, etc.")
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
		color.Println("@{!r}Se esperaba -nombre, -path, -id, etc.")
	}
}

func otroParamRep() {
	if token < (len(tokens) - 1) {
		if tokens[token+1].GetTipo() == "TK_NOMBRE" || tokens[token+1].GetTipo() == "TK_PATH" || tokens[token+1].GetTipo() == "TK_PID" || tokens[token+1].GetTipo() == "TK_RUTA" {
			paramRep()
		}
	}
}

func paramLogin() {
	tokenAux = nextToken()
	if tokenCorrecto(tokenAux, "TK_USR") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR USERNAME
				vUser = tokenAux.GetLexema()
				otroParamLogin()
			} else {
				color.Println("@{!r}El username debe ser un ID (iniciar con letra).")
				ejLogin = false
				syntaxError = true
			}
		} else {
			ejLogin = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_PWD") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_NUM") || tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR PASS
				vPass = tokenAux.GetLexema()
				otroParamLogin()
			} else {
				color.Println("@{!r}La password debe ser solo numérica o un ID (iniciar con letra).")
				ejLogin = false
				syntaxError = true
			}
		} else {
			ejLogin = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_PID") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR ID
				vID = tokenAux.GetLexema()
				otroParamLogin()
			} else {
				ejLogin = false
				syntaxError = true
			}
		} else {
			ejLogin = false
			syntaxError = true
		}
	} else {
		ejLogin = false
		syntaxError = true
		color.Println("@{!r}Se esperaba -usr, -pwd o -id")
	}
}

func otroParamLogin() {
	if token < (len(tokens) - 1) {
		if tokens[token+1].GetTipo() == "TK_USR" || tokens[token+1].GetTipo() == "TK_PWD" || tokens[token+1].GetTipo() == "TK_PID" {
			paramLogin()
		}
	}
}

func paramMkgrp() {
	tokenAux = nextToken()
	if tokenCorrecto(tokenAux, "TK_PID") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR ID
				vID = tokenAux.GetLexema()
				otroParamMkgrp()
			} else {
				ejMkgrp = false
				syntaxError = true
			}
		} else {
			ejMkgrp = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_NAME") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR ID
				vGroup = tokenAux.GetLexema()
				otroParamMkgrp()
			} else {
				color.Println("@{!r}El nombre del grupo debe ser un ID (iniciar con letra).")
				ejMkgrp = false
				syntaxError = true
			}
		} else {
			ejMkgrp = false
			syntaxError = true
		}
	} else {
		ejMkgrp = false
		syntaxError = true
		color.Println("@{!r}Se esperaba -name o -id.")
	}
}

func otroParamMkgrp() {
	if token < (len(tokens) - 1) {
		if tokens[token+1].GetTipo() == "TK_PID" || tokens[token+1].GetTipo() == "TK_NAME" {
			paramMkgrp()
		}
	}
}

func paramRmgrp() {
	tokenAux = nextToken()
	if tokenCorrecto(tokenAux, "TK_PID") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR ID
				vID = tokenAux.GetLexema()
				otroParamRmgrp()
			} else {
				ejRmgrp = false
				syntaxError = true
			}
		} else {
			ejRmgrp = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_NAME") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR ID
				vGroup = tokenAux.GetLexema()
				otroParamRmgrp()
			} else {
				color.Println("@{!r}El nombre del grupo debe ser un ID (iniciar con letra).")
				ejRmgrp = false
				syntaxError = true
			}
		} else {
			ejRmgrp = false
			syntaxError = true
		}
	} else {
		ejRmgrp = false
		syntaxError = true
		color.Println("@{!r}Se esperaba -name o -id.")
	}
}

func otroParamRmgrp() {
	if token < (len(tokens) - 1) {
		if tokens[token+1].GetTipo() == "TK_PID" || tokens[token+1].GetTipo() == "TK_NAME" {
			paramRmgrp()
		}
	}
}

func paramMkusr() {
	tokenAux = nextToken()
	if tokenCorrecto(tokenAux, "TK_PID") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR ID
				vID = tokenAux.GetLexema()
				otroParamMkusr()
			} else {
				ejMkusr = false
				syntaxError = true
			}
		} else {
			ejMkusr = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_USR") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR USERNAME
				vUser = tokenAux.GetLexema()
				otroParamMkusr()
			} else {
				color.Println("@{!r}El username debe ser un ID (iniciar con letra).")
				ejMkusr = false
				syntaxError = true
			}
		} else {
			ejMkusr = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_PWD") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_NUM") || tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR PASS
				vPass = tokenAux.GetLexema()
				otroParamMkusr()
			} else {
				color.Println("@{!r}La password debe ser solo numérica o un ID (iniciar con letra).")
				ejMkusr = false
				syntaxError = true
			}
		} else {
			ejMkusr = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_GRP") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR USERNAME
				vGroup = tokenAux.GetLexema()
				otroParamMkusr()
			} else {
				color.Println("@{!r}El grupo debe ser un ID (iniciar con letra).")
				ejMkusr = false
				syntaxError = true
			}
		} else {
			ejMkusr = false
			syntaxError = true
		}
	} else {
		ejMkusr = false
		syntaxError = true
		color.Println("@{!r}Se esperaba -id, -usr, -pwd ó -grp")
	}
}

func otroParamMkusr() {
	if token < (len(tokens) - 1) {
		if tokens[token+1].GetTipo() == "TK_PID" || tokens[token+1].GetTipo() == "TK_USR" || tokens[token+1].GetTipo() == "TK_PWD" || tokens[token+1].GetTipo() == "TK_GRP" {
			paramMkusr()
		}
	}
}

func paramRmusr() {
	tokenAux = nextToken()
	if tokenCorrecto(tokenAux, "TK_PID") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR ID
				vID = tokenAux.GetLexema()
				otroParamRmusr()
			} else {
				ejRmusr = false
				syntaxError = true
			}
		} else {
			ejRmusr = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_USR") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR USERNAME
				vUser = tokenAux.GetLexema()
				otroParamRmusr()
			} else {
				color.Println("@{!r}El username debe ser un ID (iniciar con letra).")
				ejRmusr = false
				syntaxError = true
			}
		} else {
			ejRmusr = false
			syntaxError = true
		}
	} else {
		ejRmusr = false
		syntaxError = true
		color.Println("@{!r}Se esperaba -id ó -usr")
	}
}

func otroParamRmusr() {
	if token < (len(tokens) - 1) {
		if tokens[token+1].GetTipo() == "TK_PID" || tokens[token+1].GetTipo() == "TK_USR" {
			paramRmusr()
		}
	}
}

func paramChmod() {
	tokenAux = nextToken()
	if tokenCorrecto(tokenAux, "TK_PID") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR ID
				vID = tokenAux.GetLexema()
				otroParamChmod()
			} else {
				ejChmod = false
				syntaxError = true
			}
		} else {
			ejChmod = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_PATH") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_FILE") || tokenCorrecto(tokenAux, "TK_DIR") {
				//SETEAR PATH
				vPath = tokenAux.GetLexema()
				otroParamChmod()
			} else {
				ejChmod = false
				syntaxError = true
			}
		} else {
			ejChmod = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_UGO") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_NUM") {
				//SETEAR UGO
				vUgo = tokenAux.GetLexema()
				otroParamChmod()
			} else {
				ejChmod = false
				syntaxError = true
			}
		} else {
			ejChmod = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_R") {
		//SETEAR R
		vR = tokenAux.GetLexema()
		otroParamChmod()
	} else {
		ejChmod = false
		syntaxError = true
		color.Println("@{!r}Se esperaba -id, -path, -ugo ó -r")
	}
}

func otroParamChmod() {
	if token < (len(tokens) - 1) {
		if tokens[token+1].GetTipo() == "TK_PID" || tokens[token+1].GetTipo() == "TK_PATH" || tokens[token+1].GetTipo() == "TK_UGO" || tokens[token+1].GetTipo() == "TK_R" {
			paramChmod()
		}
	}
}

func paramMkfile() {
	tokenAux = nextToken()
	if tokenCorrecto(tokenAux, "TK_PID") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR ID
				vID = tokenAux.GetLexema()
				otroParamMkfile()
			} else {
				ejMkfile = false
				syntaxError = true
			}
		} else {
			ejMkfile = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_PATH") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_FILE") {
				//SETEAR PATH
				vPath = tokenAux.GetLexema()
				otroParamMkfile()
			} else {
				ejMkfile = false
				syntaxError = true
			}
		} else {
			ejMkfile = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_SIZE") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_NUM") {
				//SETEAR SIZE
				vSize = tokenAux.GetLexema()
				otroParamMkfile()
			} else {
				ejMkfile = false
				syntaxError = true
			}
		} else {
			ejMkfile = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_CONT") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_FILE") {
				//SETEAR CONTENIDO
				vCont = tokenAux.GetLexema()
				otroParamMkfile()
			} else {
				ejMkfile = false
				syntaxError = true
			}
		} else {
			ejMkfile = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_P") {
		//SETEAR P
		vP = tokenAux.GetLexema()
		otroParamMkfile()
	} else {
		ejMkfile = false
		syntaxError = true
		color.Println("@{!r}Se esperaba -id, -path, -size, etc")
	}
}

func otroParamMkfile() {
	if token < (len(tokens) - 1) {
		if tokens[token+1].GetTipo() == "TK_PID" || tokens[token+1].GetTipo() == "TK_PATH" || tokens[token+1].GetTipo() == "TK_SIZE" || tokens[token+1].GetTipo() == "TK_CONT" || tokens[token+1].GetTipo() == "TK_P" {
			paramMkfile()
		}
	}
}

func paramRm() {
	tokenAux = nextToken()
	if tokenCorrecto(tokenAux, "TK_PID") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR ID
				vID = tokenAux.GetLexema()
				otroParamRm()
			} else {
				ejRm = false
				syntaxError = true
			}
		} else {
			ejRm = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_PATH") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_FILE") || tokenCorrecto(tokenAux, "TK_DIR") {
				//SETEAR PATH
				vPath = tokenAux.GetLexema()
				otroParamRm()
			} else {
				ejRm = false
				syntaxError = true
			}
		} else {
			ejRm = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_RF") {
		//SETEAR RF
		vRf = tokenAux.GetLexema()
		otroParamRm()
	} else {
		ejRm = false
		syntaxError = true
		color.Println("@{!r}Se esperaba -id, -path ó -rf")
	}
}

func otroParamRm() {
	if token < (len(tokens) - 1) {
		if tokens[token+1].GetTipo() == "TK_PID" || tokens[token+1].GetTipo() == "TK_PATH" || tokens[token+1].GetTipo() == "TK_RF" {
			paramRm()
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

func paramCat() {
	tokenAux = nextToken()
	if tokenCorrecto(tokenAux, "TK_PID") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR ID
				vID = tokenAux.GetLexema()
				otroParamCat()
			} else {
				ejCat = false
				syntaxError = true
			}
		} else {
			ejCat = false
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_PFILE") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_NUM") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ASIG") {
				tokenAux = nextToken()
				if tokenCorrecto(tokenAux, "TK_FILE") {
					//GUARDAR FILE
					ListaFiles = append(ListaFiles, tokenAux.GetLexema())
					otroFile()
					otroParamCat()
				} else {
					ejCat = false
					syntaxError = true
				}
			} else {
				ejCat = false
				syntaxError = true
			}
		} else {
			ejCat = false
			syntaxError = true
		}
	} else {
		ejCat = false
		syntaxError = true
		color.Println("@{!r}Se esperaba -id ó -file")
	}
}

func otroFile() {
	if token < (len(tokens) - 1) {
		if tokens[token+1].GetTipo() == "TK_PFILE" {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_PFILE") {
				tokenAux = nextToken()
				if tokenCorrecto(tokenAux, "TK_NUM") {
					tokenAux = nextToken()
					if tokenCorrecto(tokenAux, "TK_ASIG") {
						tokenAux = nextToken()
						if tokenCorrecto(tokenAux, "TK_FILE") {
							//GUARDAR FILE
							ListaFiles = append(ListaFiles, tokenAux.GetLexema())
							otroFile()
						} else {
							ejCat = false
							syntaxError = true
						}
					} else {
						ejCat = false
						syntaxError = true
					}
				} else {
					ejCat = false
					syntaxError = true
				}
			} else {
				ejCat = false
				syntaxError = true
			}
		}
	}
}

func otroParamCat() {
	if token < (len(tokens) - 1) {
		if tokens[token+1].GetTipo() == "TK_PID" || tokens[token+1].GetTipo() == "TK_PFILE" {
			paramCat()
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
		color.Println("@{!r}Se esperaba -id, -type, etc.")
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
