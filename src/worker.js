export default {
  async fetch(request, env, ctx) {
    const url = new URL(request.url);

    // API 路由处理
    if (url.pathname.startsWith("/api/")) {
      return handleAPI(request, env);
    }

    // 静态资源通过 Assets 服务
    try {
      return await env.ASSETS.fetch(request);
    } catch (e) {
      // 404 fallback 到 index.html (SPA 路由)
      if (url.pathname !== "/") {
        const indexRequest = new Request(new URL("/", request.url), request);
        return await env.ASSETS.fetch(indexRequest);
      }

      return new Response("Not Found", { status: 404 });
    }
  },
};

async function handleAPI(request, env) {
  const url = new URL(request.url);

  // 健康检查
  if (url.pathname === "/api/health") {
    return new Response(
      JSON.stringify({
        status: "ok",
        timestamp: Date.now(),
        version: "1.0.0",
      }),
      {
        headers: {
          "Content-Type": "application/json",
          "Access-Control-Allow-Origin": "*",
        },
      }
    );
  }

  // CORS 预检请求
  if (request.method === "OPTIONS") {
    return new Response(null, {
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Access-Control-Allow-Methods": "GET, POST, OPTIONS",
        "Access-Control-Allow-Headers": "Content-Type",
      },
    });
  }

  return new Response("API endpoint not found", { status: 404 });
}
