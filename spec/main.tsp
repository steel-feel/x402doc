import "@typespec/http";
import "@typespec/protobuf";

using Http;
using TypeSpec.Protobuf; 

@TypeSpec.service(#{ title: "X402Docs Services" })
@package({
      name: "x402docs.v1",
})
namespace X402Doc;

// Base model for HTTP API (includes content)
model Document {
  @field(1) id: string;
  @field(2) title: int32;
  @field(3) price : int32;
  content: string;
}

// Protobuf model (excludes content)
@message
model ProtoDocument is OmitProperties<Document, "content"> {}

model DocumentList {
   items: Document[];
}

@message
model ProtoDocumentList {
   @field(1) items: ProtoDocument[];
}

model HealthStatus {
  status: string;
}

// HTTP API Interface
@route("/api/v1")
@tag("X402Docs")
interface X402DocsHttp {
  /** List documents */
  @get list(): DocumentList;
  /** Read document */
  @get read(@path id: string): Document;
}

@route("/health")
@tag("Health")
interface HealthHttp {
  /** Check health status */
  @get check(): HealthStatus;
}

// gRPC / Protobuf Interface
@TypeSpec.Protobuf.service
interface X402DocsProto {
  /** List documents */
  list(): ProtoDocumentList;
  /** Read documents */
  read(@field(1) id: string): ProtoDocument;
}

