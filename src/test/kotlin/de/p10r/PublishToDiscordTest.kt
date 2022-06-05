package de.p10r

import io.kotest.matchers.collections.shouldHaveSize
import org.junit.jupiter.api.Test
import java.time.LocalDate

internal class PublishToDiscordTest {
    val discordServer = FakeDiscordServer()
    val discordApi = DiscordApi(
        discordApiClient = discordServer,
        today = { LocalDate.of(2022, 2, 2) }
    )

    @Test
    fun `pushes channel message to webhook uri`() {
        val schedules = Schedules(schedules = listOf(scheduleOf()))

        discordApi.publish(schedules)

        discordServer.recordedRequests shouldHaveSize 1
    }
}