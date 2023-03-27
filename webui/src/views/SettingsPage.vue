<script>
	export default {
		// Reactive State
		data(){
			return {
				ErrorMessage:null,
				Username: "",
				NewUsername :"",
				BannedUsers : []

			}
		},
		
		methods: {

			// Click sul pulsante Back
			BackButtonPressed(){
				this.$router.replace("/users/"+this.Username);
				return;
			},

			// Click sul pulsante Logout
			LogoutButtonPressed(){

				// Pulizia del localstorage
				localStorage.clear();

				// Ritorno alla schermata di Login
				this.$router.replace("/login");
				return;

			},
			
			// Click sul pulsante Confirm
			async ConfirmButtonPressed(){
				let setUsername_config =  {
					headers: {
						Authorization: `Bearer ${localStorage.getItem("Authstring")}`
					}
				};
				try{
					let setUsername_response = await this.$axios.put("/users/"+this.Username+"/username",{Username:this.NewUsername},setUsername_config);
					this.Username = setUsername_response.data.Username;
					this.NewUsername="";
					localStorage.setItem("Username",this.Username);
					this.$router.replace("/users/"+this.Username+"/settings");
					return;

				}
				catch(e){
					this.NewUsername="";
					if(e.response.status==500){
						this.ErrorMessage="Oops,esiste già un utente con questo Username,si prega di sceglierne un altro."
						return;
					}
				}


			},
			
			// Click sul pulsante Unban
			async UnbanButtonPressed(BannedId){
				let setUnbanUser_config =  {
					headers: {
						Authorization: `Bearer ${localStorage.getItem("Authstring")}`
					}
				};
				try{
					let unbanUser_response = await this.$axios.delete("/users/"+this.Username+"/bannedusers/"+BannedId,setUnbanUser_config);

					// Aggiornamento della lista BannedUsers
					this.BannedUsers = this.BannedUsers.filter(BannedUser => BannedUser != BannedId);
					return;


				}
				catch(e){
					return;
				}

			},
			
			// Click sul pulsante Delete User
			async DeleteUserButtonPressed(){
				let deleteUser_config =  {
					headers: {
						Authorization: `Bearer ${localStorage.getItem("Authstring")}`
					}
				};
				try{
					let deleteUser_response = await this.$axios.delete("/users/"+this.Username+"/",deleteUser_config);
					localStorage.clear();
					this.$router.replace("/login");
					return;

				}
				catch(e){
					return;
				}

			}
		
		},
		
		// Metodo da eseguire al montaggio del componente
		async mounted(){
			// Setto l'Username
			this.Username = this.$route.params.Username;

			// Ottengo la lista di utenti bannati

			let getBanned_config =  {
				headers: {
					Authorization: `Bearer ${localStorage.getItem("Authstring")}`
				}
			};
			try {
				let getBanned_response = await this.$axios.get("/users/"+this.Username+"/bannedusers/",getBanned_config);
				this.BannedUsers = getBanned_response.data.BannedUsers.map(x => x.BannedId);
				return;
			}
			catch(e){
				return;
			}
				
		}
	}
	
</script>

<template>
	<ErrorMsg v-if="ErrorMessage" :Message =ErrorMessage></ErrorMsg>


	<div id="SettingsPageContainer" class="container-fluid" style="display: flex; flex-direction: column;  min-width: 100vw; min-height: 100vh;">
		
		<!-- Pannello superiore: pulsanti Back e Logout -->
		<div style="display: flex; flex-direction: row; justify-content: space-between; margin-top: 2em;">
			
			<!-- Back Button -->
			<div>
				<button class="btn btn-primary" id="BackButton" @click="BackButtonPressed">
					<i class="fa solid fa-arrow-left"> Back</i>
				</button>
			</div>

			<!-- Logout Button -->
			<div>
				<button class="btn btn-dark" id="LogoutButton" @click="LogoutButtonPressed" style="background-color: darkred;">
					<i class="fa solid fa-right-from-bracket" style="color: black;"> Logout</i>
				</button>
			</div>


		</div>

		<!-- Pannello Inferiore: opzioni -->
		<div style="display: flex; flex-direction: column; align-items: center; row-gap: 50px;">
			
			<!-- Change Username -->
			<div class="card" style="display: flex; flex-direction: column; justify-content:center; align-items: center; height: 150px; width: 700px;" >
				<div style="display: flex; flex-direction: column;">
					<div>
						<label for="ChangeUsernameTextArea" class="form-label" style="font-size: medium;">Change Username</label>
					</div>
					<div style="display: flex; flex-direction: row;gap: 5px;">
						<input type="text"  v-model="NewUsername" :placeholder="'Current Username: '+Username" id="ChangeUsernameTextArea" class="form-control">
						<button  id="ChangeUsernameConfirmButton" class="btn btn-primary" @click="ConfirmButtonPressed">Confirm</button>
					</div>
					<div id="UsernameHelperBlock" class="form-text"> 
						Your username must be 3-30 characters long, contain letters and numbers, and must not contain spaces. 
					</div>
				</div>
			</div>

			<!-- View Banned Users Button-->
			<div class="card" style="display: flex; flex-direction: row; justify-content:center; align-items:center; height: 100px; width: 700px;">
				<div style="display: flex; flex-direction: row; gap:50px">
					<div style="display: flex; flex-direction: column; justify-content: center; align-items: center;">
						<button id="ViewBannedUsersButton" class="btn btn-primary" data-bs-toggle="modal" data-bs-target="#BannedUsersModal">View Banned Users</button>
					</div>
					<div style="display: flex; align-items: center;">
						<span> View and remove previously banned users</span>
					</div>				
				</div>
			</div>

			<!-- Delete User Button-->
			<div class="card" style="display: flex; flex-direction: row; justify-content: center; align-items: center; height: 50px; width: 700px;">
				<div>
					<button class="btn " id="DeleteUserButton" data-bs-toggle="modal" data-bs-target="#DeleteUserModal" style="background-color: darkred; color: white;">Delete User</button>
				</div>
			</div>

		</div>
	</div>
	
	<!-- BannedUsers Modal  -->
	<div class="modal fade" id="BannedUsersModal" tabindex="-1">
		<div class="modal-dialog modal-dialog-scrollable">
			<div class="modal-content">
				<div class="modal-header">
					<h1 class="modal-title fs-5" id="BannedUsersModalHeader">Banned Users</h1>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<div class="modal-body">
					<div v-for="BannedUser in BannedUsers" :key = BannedUser style="display: flex; flex-direction: column; align-items: center;">
						<div>
							<span>{{BannedUser}}</span>
							<button id="UnbanButton" class="btn" style="border: none;" @click="UnbanButtonPressed(BannedUser)">
								<i class="fa-regular fa-trash-can"></i>
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- DeleteUser Modal  -->
	<div class="modal fade" id="DeleteUserModal" tabindex="-1">
		<div class="modal-dialog modal-dialog-scrollable">
			<div class="modal-content">
				<div class="modal-header">
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<div class="modal-body" style="display: flex; flex-direction: column; align-items: center;">
					<span>By clicking the button below you will irreversibly delete your account.</span>
					<button class="btn" id="DeleteUserConfirmationButton" style="background-color: darkred; color: white;" @click="DeleteUserButtonPressed" >Confirm</button>
				</div>
			</div>
		</div>
	</div>
</template>

<style>
	#SettingsPageContainer{
		background: rgb(87,32,122);
		background: -moz-linear-gradient(68deg, rgba(87,32,122,1) 15%, rgba(104,20,138,1) 50%, rgba(87,32,122,1) 85%);
		background: -webkit-linear-gradient(68deg, rgba(87,32,122,1) 15%, rgba(104,20,138,1) 50%, rgba(87,32,122,1) 85%);
		background: linear-gradient(68deg, rgba(87,32,122,1) 15%, rgba(104,20,138,1) 50%, rgba(87,32,122,1) 85%);
		filter: progid:DXImageTransform.Microsoft.gradient(startColorstr="#57207a",endColorstr="#57207a",GradientType=1); 
	}

	#BackButton:hover{
		transform: scale(1.1,1.1);
	}
	#LogoutButton:hover{
		transform: scale(1.1,1.1);
	}
	#ChangeUsernameConfirmButton:hover{
		transform: scale(1.1,1.1);
	}
	#UnbanButton:hover{
		transform: scale(1.3,1.3);
		color: darkred;
	}
	#ViewBannedUsersButton:hover{
		transform: scale(1.1,1.1);
	}
	#DeleteUserButton:hover{
		transform: scale(1.1,1.1);
	}
	#DeleteUserConfirmationButton:hover{
		transform: scale(1.15,1.15);
	}
	
</style>
