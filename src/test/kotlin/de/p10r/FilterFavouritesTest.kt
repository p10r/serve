package de.p10r

import io.kotest.matchers.collections.shouldContainExactly
import org.junit.jupiter.api.Test

class FilterFavouritesTest {

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
}

