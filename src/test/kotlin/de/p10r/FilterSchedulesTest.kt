package de.p10r

import io.kotest.matchers.collections.shouldContainExactly
import io.kotest.matchers.shouldBe
import org.junit.jupiter.api.Test

class FilterSchedulesTest {

    @Test
    fun `removes all non-favoured leagues`() {
        val favourites = Schedules(
            schedules = listOf(
                Schedule("Italian Playoffs", emptyList()),
                Schedule("Russian Playoffs", emptyList()),
                Schedule("Polish Playoffs", emptyList()),
            )
        ).filterFavourites(listOf("Italian Playoffs", "Polish Playoffs"))

        favourites.shouldContainExactly(
            Schedule("Italian Playoffs", emptyList()),
            Schedule("Polish Playoffs", emptyList())
        )
    }

    @Test
    fun `removes all cancelled games`() {
        val games = scheduleOf(
            listOf(
                gameOf(isCanceled = true),
                gameOf(homeTeam = "Lube", awayTeam = "Modena")
            ),
            "league name",
        )

        games.removeCancelledGames().games shouldBe listOf(gameOf(homeTeam = "Lube", awayTeam = "Modena"))
    }
}

private fun Schedule.removeCancelledGames() = Schedule(leagueName, games.filterNot { it.isCanceled })

