package de.p10r

import org.http4k.core.Body
import org.http4k.core.HttpHandler
import org.http4k.core.Method
import org.http4k.core.Request
import org.http4k.format.Moshi.auto
import java.time.Instant

class FlashScoreApi(
    private val flashScoreClient: HttpHandler,
    private val apiKey: String,
) {
    private val bodyFrom = Body.auto<RawFlashScoreResponse>().toLens()

    fun fetchSchedule(): Schedules {
        val req = Request(Method.GET, "/v1/events/list?locale=en_GB&timezone=-4&sport_id=12&indent_days=0")
            .header("X-RapidAPI-Host", "flashscore.p.rapidapi.com")
            .header("X-RapidAPI-Key", apiKey)

        val response = flashScoreClient(req)
        println("response!")
        //response is returned
        val body = bodyFrom(response)
        //we never reach this point
        println("deserialized!")
        return body.toSchedules()
    }
}

private data class RawFlashScoreResponse(
    val DATA: List<League>,
) : List<RawFlashScoreResponse.League> by DATA {
    data class League(
        val NAME: String,
        val EVENTS: List<Event>,
    ) {
        data class Event(
            val FH: String,
            val FK: String,
            val START_TIME: Long,
        )
    }
}

private fun RawFlashScoreResponse.toSchedules() = Schedules(map { Schedule(it.NAME, it.EVENTS.toGames()) })

private fun List<RawFlashScoreResponse.League.Event>.toGames() = map {
    Game(Instant.ofEpochSecond(it.START_TIME), it.FH, it.FK)
}
