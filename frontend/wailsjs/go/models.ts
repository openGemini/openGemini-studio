export namespace main {
	
	export class AppSetting {
	    language: string;
	    theme_mode: string;
	    custom_font: string;
	    max_history_count: number;
	    data_dir: string;
	    debug: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AppSetting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.language = source["language"];
	        this.theme_mode = source["theme_mode"];
	        this.custom_font = source["custom_font"];
	        this.max_history_count = source["max_history_count"];
	        this.data_dir = source["data_dir"];
	        this.debug = source["debug"];
	    }
	}
	export class ConnectConfig {
	    name: string;
	    address: string;
	    http_schema: string;
	    ca_certificate: string;
	    client_certificate: string;
	    client_key: string;
	    insecure_tls: boolean;
	    insecure_hostname: boolean;
	    enable_auth: boolean;
	    username: string;
	    password: string;
	
	    static createFrom(source: any = {}) {
	        return new ConnectConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.address = source["address"];
	        this.http_schema = source["http_schema"];
	        this.ca_certificate = source["ca_certificate"];
	        this.client_certificate = source["client_certificate"];
	        this.client_key = source["client_key"];
	        this.insecure_tls = source["insecure_tls"];
	        this.insecure_hostname = source["insecure_hostname"];
	        this.enable_auth = source["enable_auth"];
	        this.username = source["username"];
	        this.password = source["password"];
	    }
	}
	export class RetentionPolicy {
	    name: string;
	    duration: string;
	
	    static createFrom(source: any = {}) {
	        return new RetentionPolicy(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.duration = source["duration"];
	    }
	}
	export class DatabaseMetadata {
	    retention_policies: RetentionPolicy[];
	    measurements: string[];
	
	    static createFrom(source: any = {}) {
	        return new DatabaseMetadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.retention_policies = this.convertValues(source["retention_policies"], RetentionPolicy);
	        this.measurements = source["measurements"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ExecuteRequest {
	    connect_name: string;
	    database: string;
	    retention_policy: string;
	    measurement: string;
	    precision: string;
	    command: string;
	
	    static createFrom(source: any = {}) {
	        return new ExecuteRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.connect_name = source["connect_name"];
	        this.database = source["database"];
	        this.retention_policy = source["retention_policy"];
	        this.measurement = source["measurement"];
	        this.precision = source["precision"];
	        this.command = source["command"];
	    }
	}
	export class ExecuteResponse {
	    no_content: boolean;
	    message: string;
	    execution_time: number;
	    columns: string[];
	    values: any[][];
	
	    static createFrom(source: any = {}) {
	        return new ExecuteResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.no_content = source["no_content"];
	        this.message = source["message"];
	        this.execution_time = source["execution_time"];
	        this.columns = source["columns"];
	        this.values = source["values"];
	    }
	}
	export class History {
	    id: string;
	    query: string;
	    timestamp: number;
	    execution_time: number;
	    database: string;
	    retention_policy: string;
	    success: boolean;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new History(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.query = source["query"];
	        this.timestamp = source["timestamp"];
	        this.execution_time = source["execution_time"];
	        this.database = source["database"];
	        this.retention_policy = source["retention_policy"];
	        this.success = source["success"];
	        this.error = source["error"];
	    }
	}

}

