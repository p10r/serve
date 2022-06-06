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
    val isCanceled: Boolean,
)

sealed interface ServeResult {
    object NoGamesToday : ServeResult

    data class TodaysGames(val schedules: Schedules) : ServeResult
}

fun Schedules.buildTodaysSchedule(): ServeResult = also { println("fetched schedule!") }
    .filterFavourites()
    .takeUnless { it.isEmpty() }
    ?.let { ServeResult.TodaysGames(it) }
    ?: ServeResult.NoGamesToday


fun Schedules.filterFavourites(favourites: List<String> = favouriteLeagues) =
    Schedules(filter { it.leagueName in favourites })

private val favouriteLeagues = listOf(
    "Italy: SuperLega - Play Offs",
    "Poland: PlusLiga - Play Offs",
    "France: Ligue A - Play Offs",
    "Russia: Superleague - Play Offs",
    "World: Nations League",
    "World: Nations League Women",
)
