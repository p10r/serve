package de.p10r

import io.kotest.matchers.collections.shouldNotBeEmpty
import org.junit.jupiter.api.Test

class FetchScheduleTest {
    @Test
    fun `should return today's schedule`() {
        val api = FlashScoreApi(FakeFlashScoreApi(), "apiKey")
        api.fetchSchedules().shouldNotBeEmpty()
    }
}
