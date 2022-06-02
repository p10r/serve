package de.p10r

import io.kotest.matchers.collections.shouldBeEmpty
import io.kotest.matchers.collections.shouldHaveSize
import org.http4k.core.Status
import org.http4k.kotest.shouldHaveStatus
import org.junit.jupiter.api.Test

class E2ETest {
    @Test
    fun `publishes a discord message`() {
        with(Fixture()) {
            trigger()
            discord.recordedRequests shouldHaveSize 1
        }
    }

    @Test
    fun `forwards flash score api error`() {
        with(Fixture.withOutApiKey()) {
            val response = trigger()
            discord.recordedRequests.shouldBeEmpty()
            response shouldHaveStatus Status.INTERNAL_SERVER_ERROR
        }
    }
}