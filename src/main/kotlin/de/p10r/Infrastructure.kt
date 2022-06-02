package de.p10r

import org.http4k.client.JavaHttpClient
import org.http4k.core.Uri
import org.http4k.core.then
import org.http4k.filter.ClientFilters
import org.http4k.filter.DebuggingFilters
import org.http4k.serverless.ApiGatewayV1LambdaFunction
import org.http4k.serverless.AppLoader


@Suppress("unused")
class ServeAppFunction : ApiGatewayV1LambdaFunction(AppLoader { env: Map<String, String> ->
    ProdApp(env).routes()
})

fun ProdApp(env: Map<String, String>): App {
    val flashScoreUri = env["FLASH_SCORE_URI"]
        ?: error("env variable FLASH_SCORE_URI not provided")

    val discordUri = env["DISCORD_URI"]
        ?: error("env variable DISCORD_URI not provided")

    val apiKey = env["FLASH_SCORE_API_KEY"]
        ?: error("env variable DISCORD_URI not provided")

    val flashScoreClient = ClientFilters.SetBaseUriFrom(Uri.of(flashScoreUri))
        .then(JavaHttpClient())

    val discordClient = ClientFilters.SetBaseUriFrom(Uri.of(discordUri))
        .then(DebuggingFilters.PrintRequestAndResponse())
        .then(JavaHttpClient())

    return App(discordClient, flashScoreClient, apiKey)
}