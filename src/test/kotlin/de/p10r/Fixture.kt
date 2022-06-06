package de.p10r

class Fixture(private val apiKey: String = "a-legit-api-key") {
    companion object

    val discord = FakeDiscordServer()
    private val flashScore = FakeFlashScoreApi()

    fun run() = App(discord, flashScore, apiKey)
}