package de.p10r

import io.kotest.matchers.collections.shouldHaveSize
import org.junit.jupiter.api.Test
import java.time.Instant
import java.time.LocalDate

internal class PublishToDiscordTest {
    val discordServer = FakeDiscordServer()
    val discordApi = DiscordApi(
        discordApiClient = discordServer,
        today = { LocalDate.of(2022, 2, 2) }
    )

    @Test
    fun `pushes channel message to webhook uri`() {
        val may9th = Instant.ofEpochSecond(1652104800)
        val schedules = Schedules(
            schedules = listOf(
                Schedule(
                    leagueName = "Bundesliga",
                    games = listOf(Game(time = may9th, homeTeam = "BRV", awayTeam = "VfB"))
                )
            )
        )
        
        discordApi.publish(schedules)

        discordServer.recordedRequests shouldHaveSize 1
    }
}