package de.p10r

import org.http4k.core.HttpHandler
import org.http4k.core.Request
import org.http4k.core.Response
import org.http4k.core.Status

class FakeDiscordServer: HttpHandler {
    val recordedRequests = mutableListOf<String>()

    override fun invoke(req: Request): Response {
        recordedRequests.add(req.bodyString())
        return Response(Status.OK)
    }
}