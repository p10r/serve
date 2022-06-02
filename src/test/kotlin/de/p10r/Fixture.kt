package de.p10r

import org.http4k.core.Method
import org.http4k.core.Request

class Fixture(apiKey: String = "a-legit-api-key") {
    companion object {
        fun withOutApiKey() = Fixture("")
    }

    val discord = FakeDiscordServer()
    private val flashScore = FakeFlashScoreApi()

    private val app = App(discord, flashScore, apiKey)


    fun trigger() = app.routes()(Request(Method.GET, "/trigger"))
}