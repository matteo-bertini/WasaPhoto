<script>
export default {
	// Reactive State
	data: function () {
		return {
			errormsg: null,
			loading: false,
			Username : "",
			Followers : 0,
			Following :0,
			NumberOfPhotos: 0,
			Photos : [],
			ToSearch : "",
			IsOwner : false
		}
	},
	// Dichiarazione del watcher: quando L'Username nel path cambierà,verrà caricato il profilo corrispondente 
	watch:{
        currentPath(newUsername,oldUsername){
            if (newUsername !== oldUsername){
                this.loadProfile()
            }
        },
    },

	computed:{

        currentPath(){
            return this.$route.params.Username;
        },
	},

	methods: {
		
		// Caricamento del profilo 
		async loadProfile(){

			// Controllo se chi sta visualizzando il profilo ne è il proprietario,la pagina ne viene influenzata
			if(this.$route.params.Username === `${localStorage.getItem("Username")}`){
				this.IsOwner=true;
			}
			else{
				this.IsOwner=false;
			}

			// Richiesta http
			let config = await {headers: {Authorization: `Bearer ${localStorage.getItem("Authstring")}`},params: {Username: this.$route.params.Username}}
			try{
				let response =  await this.$axios.get("/users/",config);
				this.Username=response.data.Username;
				this.Followers=response.data.Followers;
				this.Following=response.data.Following;
				this.NumberOfPhotos=response.data.NumberOfPhotos;
				this.Photos = response.data.UploadedPhotos.reverse();
				return;
			}catch(e){
				this.errormsg = e.toString();
				return;
			}
		},
		
		// Upload di una foto 
		async uploadPhoto() {
			let input_file = document.getElementById('photo_uploader').files[0];
			const reader = new FileReader();
			reader.readAsArrayBuffer(input_file);
			let config = {
				headers: {
					"Authorization": `Bearer ${localStorage.getItem("Authstring")}`,
					"Content-Type" : "image/png"
				
				
				}
			
			}
		
			reader.onload = async () => {
				try {
					let response = await this.$axios.post("/users/"+this.Username+"/photos/", reader.result,config)
					this.NumberOfPhotos +=1;
					this.Photos.unshift(response.data);
					return;


				}catch(e){
					this.errormsg = e.toString();
					return;

				}
               
				
			}
		
		},
		
		// Rimozione di una foto dalla visualizzazione
		removePhotoFromList(photoid){
			this.Photos = this.Photos.filter(photo => photo.PhotoId !== photoid);
			this.NumberOfPhotos -=1;
		},

		// Ricerca di un profilo
		async SearchProfile(){
			
			// Aggiorno l'URL e la pagina si aggiorna automaticamente con i dati giusti
			this.$router.push(this.ToSearch);
			this.ToSearch="";
		}
		
	},

	// Eseguita appena il componente è stato montato
 	async mounted()  {
		await this.loadProfile();
	},
	
	
}


</script>

<template>
	
	<div class="container-fluid" style=" display:flex; flex-direction: column; " >
		
		<!-- Titolo e barra di ricerca -->
		<div style="display: flex; flex-direction: column; align-items: center; row-gap: 10px; margin-top: 3px;">
			<span>
				WasaPhoto
			</span>
			
			<div class="input-group mb-3" style="width: fit-content;">
				<input v-model="ToSearch" type="text" class="form-control"  placeholder="Username">
				<div class="input-group-append">
				  <button @click= "SearchProfile" class="btn btn-dark" type="button">
					<i class="fa-solid fa-search"></i>
				  </button>
				</div>
			</div>
		</div>
		
		<!-- Upload Photo,Info Profilo e Settings -->
		<div style="display: flex; flex-direction: row; justify-content:space-evenly; align-items: center; margin-top: 20px;">
			
			<!-- Se chi visualizza il profilo ne è il proprietario mostro il pulsante UploadPhoto -->
			<div v-if="IsOwner">
				<input type="file" id="photo_uploader" ref="photo_uploader" @change="uploadPhoto" accept=".png" hidden/>
				<label for="photo_uploader" class="btn btn-dark" type="button" style="height: 40px; width:160px">
					<i class="fa-solid fa-upload"></i>
					Upload Photo
				</label>
			</div>
			<!-- Se chi visualizza il profilo NON ne è il proprietario mostro il pulsate Follow -->
			<div v-else>
				<button class="btn btn-dark" type="button" style="height: 40px; width:160px">
					<i class="fa-solid fa-plus"></i>
					Follow
				</button>
			</div>
			
			<div class="card" style="width: fit-content; text-align: center;">
				<div class="card-body">
					<h5 class ="card-title" style="font-size: 30px;">{{Username}}</h5>
					<p class="card-text" style="font-size: 15px;"> Post: {{NumberOfPhotos}} | Followers: {{Followers}} | Following: {{Following}} </p>
				</div>
			</div>
			
			<div>
				<button v-if="IsOwner" class="btn btn-dark" @click="GoToSetting" type="button" style="height: 40px; width:160px"> 
					<i class="fa-solid fa-gear"></i>
					Settings
				</button>
				<button v-else class="btn btn-dark" type="button" style="height: 40px; width:160px"> 
					<i class="fa-solid fa-ban"></i>
					Ban
				</button>
				
			</div>
		</div>

		<!-- Etichetta Photos -->
		<div style="display: flex; justify-content: center; margin-top: 15px;">
			<h3>Photos</h3>	
		</div>

		<!-- Collezione delle foto del profilo -->
		<div style="display: flex; flex-direction: row; flex-wrap: wrap; gap: 15px; margin-top: 20px; justify-content: center;">
			
			
			<Photo v-for="photo in Photos"
			:key = "photo.PhotoId"
			:owner = "Username"
			:photoid = "photo.PhotoId"
			:likesnumber = "photo.LikesNumber"
			:commentsnumber = "photo.CommentsNumber"
			:dateofupload = "photo.DateOfUpload"
			@photo_deleted_from_database = "removePhotoFromList"/>
		
		</div>
	
	</div>

</template>

<style>
	


</style>
