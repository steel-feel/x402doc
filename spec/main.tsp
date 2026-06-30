import "@typespec/http";
import "@typespec/protobuf";

using Http;
using TypeSpec.Protobuf; 

@TypeSpec.service(#{ title: "X402Docs Services" })
@package({
      name: "document.v1",
})
namespace X402Doc;

// Base model for HTTP API (includes content)
model Document {
  @field(1) id: string;
  @field(2) title: string;
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
  db_status: string;
  uptime: string;
}

model LoginRequest {
  username: string;
  password: string;
}

model LoginResponse {
  token: string;
}

model CreateDocumentRequest {
  title: string;
  price: int32;
  content: string;
}

model EthereumPriceResponse {
  price: float64;
}

// HTTP API Interface
@route("/api/v1")
@tag("X402Docs")
interface X402DocsHttp {
  /** List documents */
  @get list(): DocumentList;
  
  /** Read document */
  @get read(@path id: string): Document;

  /** Login with username and password to get a JWT token */
  @route("login") @post login(@body req: LoginRequest): LoginResponse;

  /** Create document */
  @route("documents") @post @useAuth(BearerAuth) create(@body doc: CreateDocumentRequest): Document;

  /** Fetch the current Ethereum price */
  @route("price/eth") @get getEthereumPrice(): EthereumPriceResponse;
}

@route("/health")
@tag("Health")
interface HealthHttp {
  /** Check health status */
  @get check(): HealthStatus;
}

model DocumentResponse {
  @field(1) id: string;
  @field(2) title: string;
  @field(3) content: string;
  @field(4) priceUsd: int32;
}

model ListDocumentsReq {}

model ListDocumentsResponse {
  @field(1) documents: DocumentResponse[];
}

// gRPC / Protobuf Interface
@route("/grpc/document")
@TypeSpec.Protobuf.service
interface DocumentService {
  @route("get") getDocument(@field(1) id: string): DocumentResponse;
  @route("list") listDocuments(@field(1) req: ListDocumentsReq): ListDocumentsResponse;
}

