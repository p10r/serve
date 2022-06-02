package de.p10r

import org.http4k.core.HttpHandler
import org.http4k.core.Method.GET
import org.http4k.core.Response
import org.http4k.core.Status.Companion.NO_CONTENT
import org.http4k.core.Status.Companion.OK
import org.http4k.core.then
import org.http4k.filter.HandleRemoteRequestFailed
import org.http4k.filter.ServerFilters
import org.http4k.filter.debug
import org.http4k.routing.RoutingHttpHandler
import org.http4k.routing.bind
import org.http4k.routing.routes
import java.time.LocalDate

class App(
    discordApi: HttpHandler,
    flashScoreApi: HttpHandler,
    apiKey: String,
) {
    private val discord = DiscordApi(discordApi, today = LocalDate::now)
    private val flashScore = FlashScoreApi(flashScoreApi.debug(), apiKey)

    fun run(): Response {
        println("starting fetch")
        return flashScore.fetchSchedule()
            .also { println("fetched schedule!") }
            .filterFavourites()
            .also { println("filtered Favourites!") }
            .takeUnless { schedules -> schedules.isEmpty() }
            ?.let { schedules -> discord.publish(schedules) }
            ?: Response(NO_CONTENT).also { println("No games today") }
    }


    fun routes(): RoutingHttpHandler = ServerFilters.CatchAll()
        .then(ServerFilters.CatchLensFailure())
        .then(ServerFilters.HandleRemoteRequestFailed())
        .then(
            routes(
                "/" bind GET to { Response(OK).body("ok") },
                "/status" bind GET to { Response(OK).body("up") },
                "/trigger" bind GET to { run() },
                "/f" bind GET to { Response(OK).body(flashScore.fetchSchedule().toString()) }
            )
        )
}
