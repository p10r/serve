package de.p10r

import org.http4k.core.HttpHandler
import org.http4k.core.Response
import org.http4k.core.Status.Companion.NO_CONTENT
import org.http4k.filter.debug
import java.time.LocalDate

fun App(
    discordApi: HttpHandler,
    flashScoreApi: HttpHandler,
    apiKey: String,
): Response {
    val discord = DiscordApi(discordApi, today = LocalDate::now)
    val flashScore = FlashScoreApi(flashScoreApi.debug(), apiKey)

    return flashScore.fetchSchedules()
        .also { println("fetched schedule!") }
        .buildTodaysSchedule()
        .let { result ->
            when (result) {
                is ServeResult.NoGamesToday -> Response(NO_CONTENT).also { println("No games today") }
                is ServeResult.TodaysGames  -> discord.publish(result.schedules)
            }
        }
}

