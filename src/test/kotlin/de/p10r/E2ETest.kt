package de.p10r

import io.kotest.assertions.throwables.shouldThrow
import io.kotest.matchers.collections.shouldBeEmpty
import io.kotest.matchers.collections.shouldHaveSize
import org.junit.jupiter.api.Test

class E2ETest {
    @Test
    fun `publishes a discord message`() {
        with(Fixture()) {
            run()
            discord.recordedRequests shouldHaveSize 1
        }
    }

    @Test
    fun `forwards flash score api error`() {
        with(Fixture.withOutApiKey()) {
            shouldThrow<RuntimeException> {
                run()
            }
            discord.recordedRequests.shouldBeEmpty()
        }
    }

    fun Fixture.Companion.withOutApiKey() = Fixture("")
}