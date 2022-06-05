package de.p10r

import java.time.Instant

data class Schedules(
    val schedules: List<Schedule>,
) : List<Schedule> by schedules

data class Schedule(
    val leagueName: String,
    val games: List<Game>,
)

data class Game(
    val time: Instant,
    val homeTeam: String,
    val awayTeam: String,
)

fun Schedules.filterFavourites(favourites: List<String> = favouriteLeagues) =
    Schedules(filter { it.leagueName in favourites })

private val favouriteLeagues = listOf(
    "Italy: SuperLega - Play Offs",
    "Poland: PlusLiga - Play Offs",
    "France: Ligue A - Play Offs",
    "Russia: Superleague - Play Offs",
    "World: Nations League",
    "World: Nations League Women",
//    "European Championships",
//    "Champions League",
//    "CEV Cup",
//    "European Games",
//    "Ligue A",
//    "1. Bundesliga",
//    "VBL Supercup",
//    "SuperLega",
//    "Super Cup",
//    "NORCECA Championship",
//    "PlusLiga",
//    "Polish Cup",
//    "Super Cup",
//    "Superleague",
//    "Superleague - Playoffs",
//    "Russia Cup",
//    "Super Cup",
//    "World Championship",
//    "World Championship Women",
//    "Olympic Games",
//    "Olympic Games Women",
//    "Nations League",
//    "Club World Championship",
)

/*
"European Championships",
"European Championships Women",
"Champions League",
"Champions League Women",
"CEV Cup",
"Ligue A",
"Supercup",
"1. Bundesliga",
"DVV Cup",
"VBL Supercup",
"SuperLega",
"Coppa Italia A1",
"Super Cup",
"Serie A1 Women",
"Coppa Italia A1 Women",
"Super Cup Women",
"NORCECA Championship",
"PlusLiga",
"Polish Cup",
"Super Cup",
"Polish Cup Women",
"Super Cup Women",
"Russia Cup",
"Super Cup",
"Efeler Ligi",
"1. Ligi",
"Turkish Cup",
"Super Cup",
"Super Cup Women",
"World Championship",
"World Championship Women",
"Olympic Games",
"Olympic Games Women",
"Nations League",
"Nations League Women",
"Club World Championship",
"Club World Championship Women",
"Hubert Wagner Memorial",
 */