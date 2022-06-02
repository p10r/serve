package de.p10r

import org.http4k.core.Body
import org.http4k.core.HttpHandler
import org.http4k.core.Method
import org.http4k.core.Request
import org.http4k.core.Response
import org.http4k.core.with
import org.http4k.format.Moshi.auto
import java.time.LocalDate

class DiscordApi(
    private val discordApiClient: HttpHandler,
    private val today: () -> LocalDate,
) {
    private val discordMessageLens = Body.auto<DiscordMessage>().toLens()

    fun publish(schedules: Schedules): Response {
        val discordMessage = discordMessageOf(schedules, today())

        return Request(Method.POST, "")
            .header("Content-Type", "application/json")
            .with(discordMessageLens of discordMessage)
            .let(discordApiClient)
    }
}

data class DiscordMessage(
    val content: String,
    val embeds: List<Embeds>,
) {
    companion object;

    data class Embeds(
        val fields: List<Fields>,
    ) {
        data class Fields(
            val name: String,
            val value: String,
            val inline: Boolean,
        )
    }
}
