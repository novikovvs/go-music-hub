export namespace main {

	export class Track {


	    static createFrom(source: any = {}) {
	        return new Track(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);

	    }
	}
	export class TrackLibrary {


	    static createFrom(source: any = {}) {
	        return new TrackLibrary(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);

	    }
	}

}






