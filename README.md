# Système de Surveillance des Capteurs Aéroportuaires (Airport MQTT)

## Description
Ce projet consiste en un système de collecte de données des capteurs météorologiques (température, vent, pression) dans un aéroport. Les données sont traitées et enregistrées en utilisant le protocole MQTT, avec une gestion des alertes et une sauvegarde des données dans une base NoSQL et des fichiers CSV.

## Architecture
Le système utilise une architecture basée sur MQTT pour la communication entre les capteurs, un gestionnaire d'alertes, un enregistreur de base de données et un enregistreur de fichiers.

## Comment lancer le projet

### Prérequis
- Go installé sur la machine.
- Installation avec une base de données InfluxDB
- Clé API Météo France

### Étapes de démarrage
1. Ajouter les clés API pour InfluxDB et Météo France dans un fichier .env **METEO_FRANCE_API_KEY** et **INFLUX_DB_API_KEY**
2. Configurez les fichiers de configuration des capteurs dans le dossier `config/` avec les paramètres souhaités.
3. Assurez-vous que le broker MQTT est opérationnel sur la machine interne.
4. Exécutez les simulateurs de capteurs pour démarrer la collecte des données :

```bash
go run ./cmd/airportsensors/pressure/pressure.go
```
```bash
go run ./cmd/airportsensors/temperature/temperature.go
```
```bash
go run ./cmd/airportsensors/wind/wind.go
```

4. Lancez le gestionnaire d'alertes :

```bash
go run ./cmd/alertmanager/alertmanager.go
```

5. Démarrez les enregistreurs de données :
```bash
go run ./cmd/databaserecorder/databaserecorder.go
```
```bash
go run ./cmd/filerecorder/filerecorder.go
```

6. Pour lancer l'API :

```
go run ./cmd/rest/rest.go
```


## Configuration
Modifiez les fichiers `.yml` dans le dossier `config/` pour ajuster les paramètres tels que l'intervalle de publication des données, les seuils d'alerte, etc.


## Membres du projet :technologist:

EGENSCHEVILLER Frédéric</br>
LABORDE Baptiste</br>
CHALEKH Zineddine</br>
GUILLOUET Adam</br>
JOLY-JEHENNE Léo
