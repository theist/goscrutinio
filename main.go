package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

type party struct {
	Color           string  `json:"color"`
	Diputados       int     `json:"diputados"`
	SiglasCortas    string  `json:"siglas_cortas"`
	Votos           int     `json:"votos"`
	PorcentajeVotos float64 `json:"porcentaje_votos"`
	Siglas          string  `json:"siglas"`
	ID              string  `json:"id"`
}

type mainJSON struct {
	VotosNulos               int     `json:"votos_nulos"`
	Code                     string  `json:"code"`
	PorcentajeCensoEscrutado float64 `json:"porcentaje_censo_escrutado"`
	PorcentajeVotosBlanco    float64 `json:"porcentaje_votos_blanco"`
	Abstencion               int     `json:"abstencion"`
	PorcentajeTotalVotantes  float64 `json:"porcentaje_total_votantes"`
	MesasTotales             int     `json:"mesas_totales"`
	TotalVotantes            int     `json:"total_votantes"`
	Year                     int     `json:"year"`
	IDRegistro               string  `json:"id_registro"`
	PorcentajeAbstencion     float64 `json:"porcentaje_abstencion"`
	Censo                    int     `json:"censo"`
	CcaaCode                 string  `json:"ccaa_code"`
	VotosBlanco              int     `json:"votos_blanco"`
	Type                     string  `json:"type"`
	CensoEscrutado           int     `json:"censo_escrutado"`
	Diputados                int     `json:"diputados"`
	PorcentajeVotosNulos     float64 `json:"porcentaje_votos_nulos"`
	Partidos                 []party `json:"partidos"`
	Month                    string  `json:"month"`
	Slug                     string  `json:"slug"`
	Name                     string  `json:"name"`
}

func getJSON(uri string) ([]byte, error) {

	resp, err := http.Get("https://elecciones.unidadeditorial.es/elecciones-generales/resultados/senado/2019/99.json")
	if err != nil {
		return []byte(""), fmt.Errorf("Error making http request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return []byte(""), fmt.Errorf("Error reponse making http request: %v", resp.Status)
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), fmt.Errorf("Error reading Body: %v", err)
	}
	return result, nil
}

func main() {

	var allParties bool

	flag.BoolVar(&allParties, "all", false, "Shows all parties even if they have no any seats")
	flag.Parse()

	res, err := getJSON("https://elecciones.unidadeditorial.es/elecciones-generales/resultados/senado/2019/99.json")
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("res\n%v", string(res))
	var jsonRes mainJSON
	err = json.Unmarshal(res, &jsonRes)
	if err != nil {
		log.Fatal("Error unmarshaling JSON", err)
	}

	fmt.Print("Resultados Generales:\n\n")
	fmt.Printf("Escrutado:    %v%%\n\n", jsonRes.PorcentajeCensoEscrutado)
	fmt.Printf("Abstencion:   %v (%v%%)\n", jsonRes.Abstencion, jsonRes.PorcentajeAbstencion)
	fmt.Printf("Voto Nulo:    %v (%v%%)\n", jsonRes.VotosNulos, jsonRes.PorcentajeVotosNulos)
	fmt.Printf("Voto Blanco:  %v (%v%%)\n\n", jsonRes.VotosBlanco, jsonRes.PorcentajeVotosBlanco)

	sort.Slice(jsonRes.Partidos[:], func(i, j int) bool {
		if jsonRes.Partidos[i].Diputados == jsonRes.Partidos[j].Diputados {
			return jsonRes.Partidos[i].Siglas < jsonRes.Partidos[j].Siglas
		}
		return jsonRes.Partidos[i].Diputados > jsonRes.Partidos[j].Diputados
	})

	for _, party := range jsonRes.Partidos {
		if !allParties && party.Diputados < 1 {
			continue
		}
		fmt.Printf("%26s: Escanos %3d, Votos: %15d(%3.2f%%)\n", party.Siglas, party.Diputados, party.Votos, party.PorcentajeVotos)
	}
}
