<script>
export default {
	// Reactive State
	data: function () {
		return {
			errormsg: null,
			loading: false,
			Username : "",
			FollowersNumber : 0,
			Followers : [],
			FollowingNumber :0,
			Following : [],
			BannedUsers : [],
			NumberOfPhotos: 0,
			Photos : [],
			ToSearch : "",
			IsOwner : false,
			IsFollowing : false,
			IsBanned : false
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
		
		// Caricamento del profilo da mostrare
		async loadProfile(){

			// Controllo se chi sta visualizzando il profilo ne è il proprietario,la pagina ne viene influenzata
			if(this.$route.params.Username === `${localStorage.getItem("Username")}`){
				this.IsOwner=true;
			}
			else{
				this.IsOwner=false;
			}
			// Richieste http

			// Impostazione del config per la richiesta getUserProfile
			let getUserProfile_config = await {headers: {Authorization: `Bearer ${localStorage.getItem("Authstring")}`},params: {Username: this.$route.params.Username}}
			try{
				
				// getUserProfile
				let getUserProfile_response =  await this.$axios.get("/users/",getUserProfile_config);
				this.Username=getUserProfile_response.data.Username;
				this.FollowersNumber=getUserProfile_response.data.Followers;
				this.FollowingNumber=getUserProfile_response.data.Following;
				this.NumberOfPhotos=getUserProfile_response.data.NumberOfPhotos;
				this.Photos = getUserProfile_response.data.UploadedPhotos.reverse();

				// getFollowers
				let getFollowers_config =  {
					headers: {
						Authorization: `Bearer ${localStorage.getItem("Authstring")}`
					}
				};
				let getFollowers_response = await this.$axios.get("/users/"+this.Username+"/followers/",getFollowers_config);
				this.Followers = getFollowers_response.data.Followers.map(x => x.FollowerId);

				//getFollowing
				let getFollowing_config =  {
					headers: {
						Authorization: `Bearer ${localStorage.getItem("Authstring")}`
					}
				};
				let getFollowing_response = await this.$axios.get("/users/"+this.Username+"/following",getFollowing_config);
				this.Following = getFollowing_response.data.Following.map(x => x.Username);

				//getBanned
				let getBanned_config =  {
					headers: {
						Authorization: `Bearer ${localStorage.getItem("Authstring")}`
					}
				};
				let getBanned_response = await this.$axios.get("/users/"+`${localStorage.getItem("Username")}`+"/bannedusers/",getBanned_config);
				this.BannedUsers = getBanned_response.data.BannedUsers.map(x => x.BannedId);

				// Si sta visualizzando il profilo di un altro utente
				if(this.IsOwner===false){

					// IsFollowing
					if(this.Followers.includes(`${localStorage.getItem("Username")}`)){
						this.IsFollowing = true;
					}
					else{
						this.IsFollowing = false;
					}

					// IsBanned
					if(this.BannedUsers.includes(this.Username)){
						this.IsBanned = true;
					}
					else{
						this.IsBanned = false;
					}
					return;
				}
				else{
					return;
				}
			}catch(e){
				console.log(e);
				return;
			}
		},
		
		// Click sul pulsante Follow/Unfollow
		async FollowUnfollowButtonPressed(){
			// Chi visualizza il profilo lo segue già
			if(this.IsFollowing==true){
				let unfollowUser_config = await {
					headers: {
						Authorization: `Bearer ${localStorage.getItem("Authstring")}`
					}
				};
				try{
					let unfollowUser_response = await this.$axios.delete("/users/"+this.Username+"/followers/"+`${localStorage.getItem("Username")}`,unfollowUser_config);
				
					// Aggiornamento del numero di followers e della lista di followers
					this.FollowersNumber --;
					this.Followers = this.Followers.filter(x=> x!==`${localStorage.getItem("Username")}`);

					// Aggiornamento del pulsante
					this.IsFollowing=false;
					return;

				}
				catch(e){
					//
				}

			}else{
				let followUser_config = await {
					headers: {
						Authorization: `Bearer ${localStorage.getItem("Authstring")}`
					}
				};
				let followUser_response = await this.$axios.post("/users/"+this.Username+"/followers/",{"FollowerId":`${localStorage.getItem("Username")}`},followUser_config);
				
				// Aggiornamento del numero di followers e della lista di followers
				this.FollowersNumber++;
				this.Followers.push(followUser_response.data.FollowerId);
				this.IsFollowing=true;
				return

			}
			

		},
		
		async BanUnbanButtonPressed(){
			// Chi visualizza il profilo lo segue già
			if(this.IsBanned==true){
				let unbanUser_config = await {
					headers: {
						Authorization: `Bearer ${localStorage.getItem("Authstring")}`
					}
				};
				try{
					let unbanUser_response = await this.$axios.delete("/users/"+`${localStorage.getItem("Username")}`+"/bannedusers/"+this.Username,unbanUser_config);
				
					// Aggiornamento della lista dei bannati
					this.BannedUsers = this.BannedUsers.filter(x=> x!== this.Username);

					// Aggiornamento del pulsante
					this.IsBanned=false;
					return;

				}
				catch(e){
					//
				}

			}else{
				let banUser_config = await {
					headers: {
						Authorization: `Bearer ${localStorage.getItem("Authstring")}`
					}
				};
				let banUser_response = await this.$axios.post("/users/"+`${localStorage.getItem("Username")}`+"/bannedusers/",{"BannedId":this.Username},banUser_config);
				
				// Aggiornamento della lista dei bannati
				this.BannedUsers.push(banUser_response.data.BannedId);
				this.IsBanned=true;
				return

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
				<button class="btn btn-dark" type="button" id="FollowUnfollowButton" @click="FollowUnfollowButtonPressed" style="height: 40px; width:160px">
					<i class="fa-solid fa-plus" v-if = "IsFollowing == false" id="FollowIcon"> Follow</i>
					<i class="fa-solid fa-minus" v-else id="UnfollowIcon"> Unfollow</i>
				</button>
			</div>
			
			<div class="card" style="width: fit-content; text-align: center;">
				<div class="card-body">
					<h5 class ="card-title" style="font-size: 30px;">{{Username}}</h5>
					<p class="card-text" style="font-size: 15px;">
						Post: {{NumberOfPhotos}} | 
						<button data-bs-toggle="modal" :data-bs-target="'#FollowersModal' +Username" style="border:none; background-color: rgba(255, 255, 255, 0);">
							Followers: {{FollowersNumber}} 
						</button>
						<button data-bs-toggle="modal" :data-bs-target="'#FollowingModal' +Username" style="border:none; background-color: rgba(255, 255, 255, 0);">
							| Following: {{FollowingNumber}} 
						</button>
						
					</p>
				</div>
			</div>
			
			<div>
				<button v-if="IsOwner" class="btn btn-dark" @click="GoToSettings" type="button" style="height: 40px; width:160px"> 
					<i class="fa-solid fa-gear"></i>
					Settings
				</button>
				<button v-else class="btn btn-dark" @click="BanUnbanButtonPressed" type="button" style="height:40px; width:160px;"> 
					<i class="fa-solid fa-ban" v-if="IsBanned==false" style="color: darkred;"> Ban</i>
					<i class="fa-solid fa-check" v-else style="color: darkgreen;"> Unban</i>
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

		<!-- Followers Modal  -->
		<div class="modal fade" :id="'FollowersModal'+Username" tabindex="-1">
			<div class="modal-dialog modal-dialog-scrollable">
				<div class="modal-content">
					<div class="modal-header">
						<h1 class="modal-title fs-5" id="exampleModalLabel">Followers</h1>
						<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
					</div>
					<div class="modal-body">
						{{Followers}}
					</div>
				</div>
			</div>
		</div>

		<!-- Following Modal  -->
		<div class="modal fade" :id="'FollowingModal'+Username" tabindex="-1">
			<div class="modal-dialog modal-dialog-scrollable">
				<div class="modal-content">
					<div class="modal-header">
						<h1 class="modal-title fs-5" id="exampleModalLabel">Following</h1>
						<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
					</div>
					<div class="modal-body">
						{{Following}}
					</div>
				</div>
			</div>
		</div>

	
	</div>

</template>

<style>
	


</style>
