export namespace main {
	
	export class FileInfo {
	    path: string;
	    content: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.content = source["content"];
	        this.name = source["name"];
	    }
	}

}

