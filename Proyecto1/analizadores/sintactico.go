package analizadores

import (
	"Proyecto1/estructuras"
	"fmt"
)

var (
	syntaxError                                                 bool = false
	pSize, pPath, pName, pUnit, pType, pFit, pDelete, pAdd, pID bool = false, false, false, false, false, false, false, false, false
	token                                                       int  = -1
	tokenAux                                                    *estructuras.Token
)

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
		fmt.Println("Analisis sintáctico exitoso")
	} else {
		fmt.Println("Error sintáctico encontrado")
	}

}

func inicio() {

	if tokenAux.GetTipo() == "TK_CMT" {
		fmt.Println(tokenAux.GetLexema())
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
		paramMkDisk()
		otraInstruccion()
	} else if tokenAux.GetTipo() == "TK_FDISK" {
		paramFDisk()
		otraInstruccion()
	} else if tokenAux.GetTipo() == "TK_MNT" {
		paramMount()
		otraInstruccion()
	} else if tokenAux.GetTipo() == "TK_UMNT" {

		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_PID") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_NUM") {
				tokenAux = nextToken()
				if tokenCorrecto(tokenAux, "TK_ASIG") {
					tokenAux = nextToken()
					if tokenCorrecto(tokenAux, "TK_ID") {
						//Guardar ID
						otroID()
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
				otroParamMkDisk()
			} else {
				syntaxError = true
			}
		} else {
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_PATH") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_DIR") {
				//SETEAR PATH
				otroParamMkDisk()
			} else {
				syntaxError = true
			}
		} else {
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_NAME") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR NAME
				otroParamMkDisk()
			} else {
				syntaxError = true
			}
		} else {
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_UNIT") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_BYTES") {
				//SETEAR BYTES
				otroParamMkDisk()
			} else {
				syntaxError = true
			}
		} else {
			syntaxError = true
		}
	} else {
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
				otroParamFDisk()
			} else {
				syntaxError = true
			}
		} else {
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_UNIT") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_BYTES") {
				//SETEAR BYTES
				otroParamFDisk()
			} else {
				syntaxError = true
			}
		} else {
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_PATH") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_FILE") {
				//SETEAR PATH
				otroParamFDisk()
			} else {
				syntaxError = true
			}
		} else {
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_TYPE") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_PEL") {
				//SETEAR TYPE
				otroParamFDisk()
			} else {
				syntaxError = true
			}
		} else {
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_FIT") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_BFW") {
				//SETEAR FIT
				otroParamFDisk()
			} else {
				syntaxError = true
			}
		} else {
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_DEL") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_FF") {
				//SETEAR DELETE MODE
				otroParamFDisk()
			} else {
				syntaxError = true
			}
		} else {
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_NAME") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR NAME
				otroParamFDisk()
			} else {
				syntaxError = true
			}
		} else {
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_ADD") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_NUM") {
				//SETEAR NUM
				otroParamFDisk()
			} else {
				syntaxError = true
			}
		} else {
			syntaxError = true
		}
	} else {
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
				otroParamMount()
			} else {
				syntaxError = true
			}
		} else {
			syntaxError = true
		}
	} else if tokenCorrecto(tokenAux, "TK_NAME") {
		tokenAux = nextToken()
		if tokenCorrecto(tokenAux, "TK_ASIG") {
			tokenAux = nextToken()
			if tokenCorrecto(tokenAux, "TK_ID") {
				//SETEAR NAME
				otroParamMount()
			} else {
				syntaxError = true
			}
		} else {
			syntaxError = true
		}
	} else {
		syntaxError = true
		fmt.Println("Se esperaba -path o -name.")
	}
}

func otroParamMount() {
	if token < (len(tokens) - 1) {
		if tokens[token+1].GetTipo() == "TK_PATH" || tokens[token+1].GetTipo() == "TK_NAME" {
			paramMkDisk()
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
							otroID()
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
