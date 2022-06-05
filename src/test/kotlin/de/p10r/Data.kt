package de.p10r

import java.time.Instant


fun scheduleOf(
    games: List<Game> = listOf(gameOf()),
    leagueName: String = "Some League",
) = Schedule(leagueName, games)

val may9th22 = Instant.ofEpochSecond(1652104800)

fun gameOf(
    time: Instant = may9th22,
    homeTeam: String = "Kazan",
    awayTeam: String = "Moscow",
    isCanceled: Boolean = false,
) = Game(time, homeTeam, awayTeam, isCanceled)